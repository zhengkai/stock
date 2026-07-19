package web

import (
	"fmt"
	"net/http"
	"server/config"
	"time"
)

func Server() {

	mux := http.NewServeMux()

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
