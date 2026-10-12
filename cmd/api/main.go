package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-cost-optimizer/internal/config"
	"ai-cost-optimizer/internal/routing"
	"ai-cost-optimizer/internal/security"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("system boot failed: %v", err)
	}

	mux := http.NewServeMux()

	routerHandler := routing.NewRouteHandler(cfg.EncryptionKey)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("POST /v1/chat/completions", routerHandler.HandleInference)

	secureHandler := security.SecurityHeadersHandler(mux)

	server := &http.Server{
		Addr:         ":" + cfg.APIPort,
		Handler:      secureHandler,
		TLSConfig:    config.GetSecureTLSConfig(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if cfg.TLSCertPath != "" && cfg.TLSKeyPath != "" {
			if err := server.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("secure server execution error: %v", err)
			}
		} else {
			if cfg.Environment == "production" {
				log.Fatal("production requires valid TLS certificates")
			}
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server execution error: %v", err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced server shutdown executed: %v", err)
	}
}
