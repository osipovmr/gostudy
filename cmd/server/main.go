package main

import (
	"context"
	"gostudy/internal/repository"
	"gostudy/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gostudy/internal/handler"
)

func gracefulShutdown(ctx context.Context, srv *http.Server) {
	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}

func main() {

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system env")
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	userRepo := repository.NewUserRepository()
	userSvc := service.NewUserService(userRepo)

	router.GET("/api/v1/health", handler.HealthCheck)
	router.GET("/api/v1/time", handler.CurrentTime)

	handler.RegisterUserRoutes(router, userSvc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			stop()
		}
	}()

	gracefulShutdown(ctx, server)
}
