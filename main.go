package main

import (
	"balancer/proxy"
	"log"
	"net/http"
	"strconv"
)

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}
	err = config.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}
	config.Print()
	sem = make(chan struct{}, config.MaxAllowed) // 限流用

	// TODO 注册路由
	mux := http.NewServeMux()
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}

		// start healthCheck
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}

		mux.Handle(l.Pattern, MaxAllowedMiddleware(httpProxy))
	}

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: mux,
	}

	if config.Schema == "http" {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if config.Schema == "https" {
		err := srv.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}
