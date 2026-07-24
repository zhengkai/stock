package web

import (
	"crypto/sha256"
	"io"
	"net/http"
	"project/metrics"
	"project/pb"
	"project/util"
	"project/zj"

	"google.golang.org/protobuf/proto"
)

func apiTest(w http.ResponseWriter, r *http.Request) {

	metrics.ReqConcurrentInc()
	defer metrics.ReqConcurrentDec()

	util.WebPushAll(`hello`, `world2`)
}

func apiSub(w http.ResponseWriter, r *http.Request) {
	metrics.ReqConcurrentInc()
	defer metrics.ReqConcurrentDec()

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

	zj.J(`api sub`)
	zj.J(util.JSON(d))
}
