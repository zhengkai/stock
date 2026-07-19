package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"project/config"
	"project/pb"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/golang-jwt/jwt/v5"
)

func WebPush(d *pb.VAPIDSubscription, w io.Writer) {

	payload := []byte(`{"title": "Hello", "body": "world"}`)

	s := &webpush.Subscription{
		Endpoint: d.GetEndpoint(),
		Keys: webpush.Keys{
			P256dh: d.GetP256Dh(),
			Auth:   d.GetAuth(),
		},
	}

	// Send Notification
	resp, err := webpush.SendNotification(payload, s, &webpush.Options{
		Subscriber:      "zhengkai@gmail.com",
		VAPIDPublicKey:  config.VapidPublicKey,
		VAPIDPrivateKey: config.VapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		fmt.Println(`webpush.SendNotification error`, err)
		return
	}
	defer resp.Body.Close()
}

func WebPush2(d *pb.VAPIDSubscription, w io.Writer) {

	payload := []byte(`{"title": "Hello", "body": "world"}`)
	ep, err := encryptPayload(payload, d.GetP256Dh(), d.GetAuth())
	if err != nil {
		panic(err)
	}

	fmt.Println(`encrypted payload`, len(ep))

	pKey, err := privateKeyFromBase64(config.VapidPrivateKey)
	if err != nil {
		panic(err)
	}

	jwt, err := createVAPIDJWT(pKey, d.GetEndpoint())
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		d.GetEndpoint(),
		bytes.NewReader(ep),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set(`TTL`, `60`)
	req.Header.Set(`Content-Encoding`, `aes128gcm`)
	req.Header.Set(`Content-Type`, `application/octet-stream`)
	req.Header.Set(`Authorization`, `vapid t=`+jwt+`, k=`+config.VapidPublicKey)

	fmt.Println(`http req`)

	rsp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	ab, err := io.ReadAll(rsp.Body)

	fmt.Fprintf(w, "Status: %d\n", rsp.StatusCode)
	fmt.Fprintf(w, "Body: %d\n", len(ab))
	w.Write(ab)

	fmt.Println(rsp.StatusCode)
	fmt.Println(len(ab), string(ab), err)
}

func encryptPayload(payload []byte, p256dh string, auth string) ([]byte, error) {
	// 1. 解码客户端公钥（未压缩点，65字节）
	uaPublic, err := base64.RawURLEncoding.DecodeString(p256dh)
	if err != nil {
		return nil, err
	}

	// 2. 解码 auth 秘密
	authSecret, err := base64.RawURLEncoding.DecodeString(auth)
	if err != nil {
		return nil, err
	}

	// 3. 生成服务端临时 ECDH 密钥对
	curve := ecdh.P256()
	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	serverPub := serverPriv.PublicKey()
	asPublic := serverPub.Bytes() // 未压缩点，65字节

	// 4. 计算 ECDH 共享秘密
	clientPub, err := curve.NewPublicKey(uaPublic)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		return nil, err
	}

	// 5. 密钥派生（RFC 8291 Section 3.4）
	// PRK = HKDF-Extract(salt=auth_secret, IKM=shared_secret)
	prk, err := hkdf.Extract(sha256.New, sharedSecret, authSecret)
	if err != nil {
		return nil, err
	}

	// key_info = "WebPush: info" || 0x00 || ua_public || as_public
	keyInfo := append([]byte("WebPush: info\x00"), uaPublic...)
	keyInfo = append(keyInfo, asPublic...)

	// nonce_info = "WebPush: nonce" || 0x00 || ua_public || as_public
	nonceInfo := append([]byte("WebPush: nonce\x00"), uaPublic...)
	nonceInfo = append(nonceInfo, asPublic...)

	// CEK = HKDF-Expand(PRK, key_info, 16)
	cek, err := hkdf.Expand(sha256.New, prk, string(keyInfo), 16)
	if err != nil {
		return nil, err
	}

	// NONCE = HKDF-Expand(PRK, nonce_info, 12)
	nonce, err := hkdf.Expand(sha256.New, prk, string(nonceInfo), 12)
	if err != nil {
		return nil, err
	}

	// 6. 生成随机盐（仅写入记录，不参与加密）
	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return nil, err
	}

	// 7. AES-128-GCM 加密
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 明文末尾追加分隔符 0x02（RFC 8188 最终记录标记）
	plaintext := append(payload, 0x02)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// 8. 组装 aes128gcm 消息体
	// 格式：salt(16) | rs(4) | idlen(1) | keyid | ciphertext
	keyID := asPublic // 服务端公钥作为 keyid
	if len(keyID) > 255 {
		return nil, errors.New("public key too large")
	}

	var body bytes.Buffer
	body.Write(salt)

	// record size 固定 4096
	rs := make([]byte, 4)
	binary.BigEndian.PutUint32(rs, 4096)
	body.Write(rs)

	body.WriteByte(byte(len(keyID)))
	body.Write(keyID)
	body.Write(ciphertext)

	return body.Bytes(), nil
}

func createVAPIDJWT(privateKey *ecdsa.PrivateKey, endpoint string) (string, error) {

	audience, err := vapidAudience(endpoint)
	if err != nil {
		return ``, err
	}

	fmt.Println(`audience`, audience)

	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.MapClaims{
			"aud": audience,
			"exp": time.Now().
				Add(12 * time.Hour).
				Unix(),
			"sub": "mailto:admin@example.com",
		},
	)

	return token.SignedString(privateKey)
}

func vapidAudience(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil
}

func encryptPayloadX(payload []byte, p256dh string, auth string) ([]byte, error) {

	// client public key
	p256dhAB, err := base64.RawURLEncoding.DecodeString(p256dh)
	if err != nil {
		return nil, err
	}

	// auth secret
	authAB, err := base64.RawURLEncoding.DecodeString(auth)
	if err != nil {
		return nil, err
	}

	// 1. server ephemeral ECDH key
	curve := ecdh.P256()

	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	serverPub := serverPriv.PublicKey()

	// 2. ECDH shared secret
	clientPub, err := curve.NewPublicKey(p256dhAB)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		return nil, err
	}

	// 3. HKDF auth secret

	ikm, err := hkdf.Key(
		sha256.New,
		sharedSecret,
		authAB,
		"Content-Encoding: auth\x00",
		32,
	)
	if err != nil {
		return nil, err
	}

	// 4. salt

	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return nil, err
	}

	// 5. derive content encryption key

	key, err := hkdf.Key(
		sha256.New,
		ikm,
		salt,
		"Content-Encoding: aes128gcm\x00",
		16,
	)
	if err != nil {
		return nil, err
	}

	// 6. nonce

	nonce, err := hkdf.Key(
		sha256.New,
		ikm,
		salt,
		"Content-Encoding: nonce\x00",
		12,
	)
	if err != nil {
		return nil, err
	}

	// 7. AES-128-GCM

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// RFC8291 record:
	// payload + delimiter 0x02

	plaintext := append(payload, 0x02)

	ciphertext := gcm.Seal(
		nil,
		nonce,
		plaintext,
		nil,
	)

	// 8. aes128gcm body format:
	//
	// salt(16)
	// rs(4)
	// idlen(1)
	// keyid
	// ciphertext

	keyID := serverPub.Bytes()

	if len(keyID) > 255 {
		return nil, errors.New("public key too large")
	}

	var body bytes.Buffer

	// salt: 16 bytes
	body.Write(salt)

	// rs: 4 bytes
	var rsBuf [4]byte
	binary.BigEndian.PutUint32(rsBuf[:], 4096)
	body.Write(rsBuf[:])

	// keyid length: 1 byte
	body.WriteByte(byte(len(keyID)))

	// keyid
	body.Write(keyID)

	// ciphertext
	body.Write(ciphertext)

	return body.Bytes(), nil
}

func privateKeyFromBase64(s string) (*ecdsa.PrivateKey, error) {

	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(raw) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(raw))
	}

	// 用新的 ecdh API 计算公钥
	ecdhPrivate, err := ecdh.P256().NewPrivateKey(raw)
	if err != nil {
		return nil, err
	}

	pub := ecdhPrivate.PublicKey().Bytes()

	// P-256 uncompressed point:
	// 04 || X(32 bytes) || Y(32 bytes)
	if len(pub) != 65 || pub[0] != 4 {
		return nil, fmt.Errorf("invalid public key")
	}

	x := new(big.Int).SetBytes(pub[1:33])
	y := new(big.Int).SetBytes(pub[33:65])

	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(raw),
	}, nil
}
