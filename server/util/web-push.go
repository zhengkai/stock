package util

import (
	"project/config"
	"project/pb"
	"project/zj"
	"strings"

	"github.com/zhengkai/webpush-go"
)

var webpushOptions = webpush.Options{
	Subscriber:      "zhengkai@gmail.com",
	VAPIDPublicKey:  config.VapidPublicKey,
	VAPIDPrivateKey: config.VapidPrivateKey,
	TTL:             30,
}

var staticSub = NewFile(`sub`)

func WebPushAll(title, body string) {
	fl, err := staticSub.ReadDir(func(path string) bool {
		return strings.HasSuffix(path, `.pb`)
	})
	if err != nil {
		return
	}

	msg := pb.WebPush_builder{
		Title: &title,
		Body:  &body,
	}.Build()

	for _, v := range fl {
		f := NewFile(v)
		d := &pb.VAPIDSubscription{}
		err := f.ReadProto(d)
		if err != nil {
			continue
		}
		zj.J(`web push`, f)
		WebPush(JSONBin(msg), d)
	}
}

func WebPush(payload []byte, d *pb.VAPIDSubscription) {

	s := &webpush.Subscription{
		Endpoint: d.GetEndpoint(),
		Keys: webpush.Keys{
			P256dh: d.GetP256Dh(),
			Auth:   d.GetAuth(),
		},
	}

	rsp, err := webpush.SendNotification(payload, s, &webpushOptions)
	if err != nil {
		zj.W(`webpush.SendNotification error`, err)
		return
	}
	defer rsp.Body.Close()
}
