package util

import (
	"fmt"
	"project/config"
	"project/pb"
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
	for _, v := range fl {
		f := NewFile(v)
		d := &pb.VAPIDSubscription{}
		err := f.ReadProto(d)
		if err != nil {
			continue
		}
		fmt.Println(`web push`, f)
		WebPush(d)
	}
}

func WebPush(d *pb.VAPIDSubscription) {

	payload := []byte(`{"title": "Hello", "body": "world"}`)

	s := &webpush.Subscription{
		Endpoint: d.GetEndpoint(),
		Keys: webpush.Keys{
			P256dh: d.GetP256Dh(),
			Auth:   d.GetAuth(),
		},
	}

	rsp, err := webpush.SendNotification(payload, s, &webpushOptions)
	if err != nil {
		fmt.Println(`webpush.SendNotification error`, err)
		return
	}
	defer rsp.Body.Close()
}
