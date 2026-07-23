// Package web 服务

package web

import (
	"fmt"
	"net/http"
	"project/config"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Server() {

	mux := http.NewServeMux()

	mux.Handle(`/_metrics`, promhttp.Handler())
	mux.HandleFunc(`/api/test`, apiTest)
	mux.HandleFunc(`/api/sub`, apiSub)

	s := &http.Server{
		Addr:         config.Web,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	fmt.Println(`start web server`, s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
