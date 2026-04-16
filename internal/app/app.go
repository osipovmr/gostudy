package app

import (
	"gostudy/internal/config"
	"gostudy/internal/db"
	"gostudy/internal/handler"
	"gostudy/internal/repository"
	"gostudy/internal/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type App struct {
	Server *http.Server
	db     *pgxpool.Pool
}

func New(cfg *config.Config) *App {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	pool, err := db.NewPool(cfg.DBURL)
	if err != nil {
		panic(err)
	}
	db.RunMigrations(cfg.DBURL)
	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)

	router.GET("/api/v1/health", handler.HealthCheck)
	router.GET("/api/v1/time", handler.CurrentTime)
	handler.RegisterUserRoutes(router, userSvc)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &App{Server: server}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		slog.Info("starting server", "addr", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()

	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	defer a.db.Close()
	if err := a.Server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	slog.Info("server stopped")
	return nil
}
