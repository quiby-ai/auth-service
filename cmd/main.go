package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Auth Service starting on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}
