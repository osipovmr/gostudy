package app

import (
	"gostudy/internal/config"
	"gostudy/internal/db"
	"gostudy/internal/facade"
	"gostudy/internal/handler"
	"gostudy/internal/middleware"
	"gostudy/internal/repository"
	"gostudy/internal/service"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type App struct {
	server *http.Server
	db     *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	pool, err := db.NewPool(cfg.DBURL)
	if err != nil {
		return nil, err
	}
	txManager := db.NewTxManager(pool)

	if err := db.RunMigrations(cfg.DBURL); err != nil {
		slog.Error("migrations failed", "err", err)
		return nil, err
	}

	// --- DI ---
	userRepository := repository.NewUserRepository(pool)
	tokenRepository := repository.NewTokenRepository(pool)
	userService := service.NewUserService(userRepository, txManager)
	tokenService := service.NewTokenService(cfg.AccessSecret, cfg.RefreshSecret, cfg.AccessTTL, cfg.RefreshTTL, tokenRepository)
	authFacade := facade.NewAuthFacade(userService, tokenService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authFacade)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// --- Router ---
	router = setupRouter(pool, userHandler, authHandler, authMiddleware.Handler())

	srv := &http.Server{
		// Адрес для прослушивания в формате "host:port" (например, ":8080" или "0.0.0.0:8080")
		Addr: cfg.HTTPAddr,
		// HTTP роутер/хендлер, который обрабатывает все входящие запросы
		Handler: router,
		// Максимальное время ожидания чтения полного HTTP запроса от клиента
		// Если клиент не отправит запрос за 5 секунд - соединение закрывается
		ReadTimeout: cfg.ReadTimeout,
		// Максимальное время для отправки HTTP ответа клиенту после начала обработки
		// Защищает от "долгих" ответов, которые клиент может не дождаться
		WriteTimeout: cfg.WriteTimeout,
		// Максимальное время простоя TCP соединения между запросами от одного клиента
		// Если клиент не отправляет новый запрос 120 секунд - соединение закрывается
		// Важно для освобождения ресурсов при keep-alive соединениях
		IdleTimeout: cfg.IdleTimeout,
	}

	return &App{
		server: srv,
		db:     pool,
	}, nil
}

func (a *App) Run(ctx context.Context) error {

	errCh := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}
	return a.shutdown()

}

func (a *App) shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	slog.Info("shutting down server")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	a.db.Close()
	slog.Info("server stopped")
	return nil

}

func setupRouter(
	pool *pgxpool.Pool,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db not ready"})
			return
		}
		c.Status(http.StatusOK)
	})

	api := router.Group("/api/v1")
	{
		api.GET("/time", handler.CurrentTime)

		auth := api.Group("")
		auth.Use(authMiddleware)
		{
			userHandler.RegisterRoutes(auth)
		}

		authHandler.RegisterRoutes(api, authMiddleware)
	}

	return router
}
