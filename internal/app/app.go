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
	Db     *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	pool, err := db.NewPool(cfg.DBURL)
	if err != nil {
		return nil, err
	}

	if err := db.RunMigrations(cfg.DBURL); err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	api := router.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)
		api.GET("/time", handler.CurrentTime)

		userHandler.RegisterRoutes(router)
	}

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	return &App{
		Server: srv,
		Db:     pool,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	slog.Info("shutting down server")
	if err := a.Server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	a.Db.Close()
	slog.Info("server stopped")
	return nil
}
