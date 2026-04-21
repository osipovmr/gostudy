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
		slog.Error("migrations failed", "err", err)
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
		// Адрес для прослушивания в формате "host:port" (например, ":8080" или "0.0.0.0:8080")
		Addr: cfg.HTTPAddr,
		// HTTP роутер/хендлер, который обрабатывает все входящие запросы
		Handler: router,
		// Максимальное время ожидания чтения полного HTTP запроса от клиента
		// Если клиент не отправит запрос за 5 секунд - соединение закрывается
		ReadTimeout: 5 * time.Second,
		// Максимальное время для отправки HTTP ответа клиенту после начала обработки
		// Защищает от "долгих" ответов, которые клиент может не дождаться
		WriteTimeout: 10 * time.Second,
		// Максимальное время простоя TCP соединения между запросами от одного клиента
		// Если клиент не отправляет новый запрос 120 секунд - соединение закрывается
		// Важно для освобождения ресурсов при keep-alive соединениях
		IdleTimeout: 120 * time.Second,
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
