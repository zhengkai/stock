package web

import (
	"fmt"
	"io"
	"net/http"
	"project/pb"
	"project/util"

	"google.golang.org/protobuf/proto"
)

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

	fmt.Println(`api sub`)
	fmt.Println(util.JSON(d))
}
