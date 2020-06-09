package monitor

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMonitor() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: promhttp.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	}()
}
