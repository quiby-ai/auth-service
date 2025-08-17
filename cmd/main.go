package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/quiby-ai/auth-service/internal/config"
	"github.com/quiby-ai/auth-service/internal/database"
	"github.com/quiby-ai/auth-service/internal/handler"
	"github.com/quiby-ai/auth-service/internal/models"
	"github.com/quiby-ai/common/pkg/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	db, err := database.NewPostgresConnection(context.Background(), cfg.PGDSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.ClosePostgresConnection(db)

	userRepo := models.NewUserRepository(db)

	jwtConfig := &auth.JWTConfig{
		Issuer:    cfg.JWTIssuer,
		Audience:  cfg.JWTAudience,
		AccessTTL: cfg.JWTAccessTTL,
		SecretKey: cfg.JWTSecret,
	}

	r := chi.NewRouter()

	loginHandler := auth.TelegramAuthMiddleware(cfg.TelegramBotToken)(
		handler.LoginWithTelegram(jwtConfig, userRepo),
	)

	meHandler := auth.RequireAuth(jwtConfig, handler.Me(userRepo))

	r.Post("/", loginHandler.ServeHTTP)
	r.Get("/me", meHandler.ServeHTTP)
	r.Get("/healthz", handler.HealthHandler)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("Exited cleanly")
}
