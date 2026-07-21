package web

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"project/pb"
	"project/util"

	"google.golang.org/protobuf/proto"
)

func apiTest(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(`start web push`)
	// util.WebPushAll(`hello`, `world2`)
	// fmt.Println(`end web push`)
}

func apiSub(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ab, err := io.ReadAll(io.LimitReader(r.Body, 1e4)) // Limit to 10KB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(ab) < 100 {
		http.Error(w, `too short`, http.StatusBadRequest)
		return
	}

	d := &pb.VAPIDSubscription{}

	err = proto.Unmarshal(ab, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256(ab)
	f := util.NewFileF(`sub/%x.pb`, hash[:6])
	f.Write(ab)

	f = util.NewFileF(`sub/%x.json`, hash[:6])
	f.WriteJSON(d)

	fmt.Println(`api sub`)
	fmt.Println(util.JSON(d))
}
