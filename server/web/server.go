// Package web 服务

package web

import (
	"net/http"
	"project/app"
	"project/config"
	"project/util"
	"project/zj"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Server() {

	mux := http.NewServeMux()

	mux.Handle(`/_metrics`, promhttp.Handler())
	mux.HandleFunc(`/api/pool`, apiPool)
	mux.HandleFunc(`/api/sub`, apiSub)

	s := &http.Server{
		Addr:         config.Web,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	zj.J(`start web server`, s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		zj.W(err)
		return
	}
}

func apiPool(w http.ResponseWriter, r *http.Request) {
	d := app.GetPB()

	util.HTTPReturnProto(d, w, r)
}
