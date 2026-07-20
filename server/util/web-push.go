package util

import (
	"fmt"
	"io"
	"project/config"
	"project/pb"

	"github.com/zhengkai/webpush-go"
)

var webpushOptions = webpush.Options{
	Subscriber:      "zhengkai@gmail.com",
	VAPIDPublicKey:  config.VapidPublicKey,
	VAPIDPrivateKey: config.VapidPrivateKey,
	TTL:             30,
}

func WebPush(d *pb.VAPIDSubscription, w io.Writer) {

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
