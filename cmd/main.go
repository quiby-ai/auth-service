package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clientpulse-org/auth-service/internal/config"
	"github.com/clientpulse-org/auth-service/internal/handler"
	"github.com/clientpulse-org/common/pkg/auth"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	jwtConfig := auth.DefaultJWTConfig(cfg.JWTSecret)

	r := mux.NewRouter()

	r.HandleFunc("/healthz", handler.HealthHandler).Methods("GET")

	r.Handle("/", auth.TelegramAuthMiddleware(cfg.BotToken)(
		handler.LoginWithTelegram(jwtConfig),
	)).Methods("POST")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Auth Service starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down…")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("Exited cleanly")
}
