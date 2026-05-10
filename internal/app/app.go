package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/osipovmr/gostudy/internal/config"
	"github.com/osipovmr/gostudy/internal/db"
	"github.com/osipovmr/gostudy/internal/facade"
	"github.com/osipovmr/gostudy/internal/handler"
	"github.com/osipovmr/gostudy/internal/kafka"
	"github.com/osipovmr/gostudy/internal/middleware"
	"github.com/osipovmr/gostudy/internal/repository"
	"github.com/osipovmr/gostudy/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type App struct {
	server       *http.Server
	db           *pgxpool.Pool
	mailConsumer *kafka.MailConsumer
}

func New(cfg *config.Config) (*App, error) {
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
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	var kafkaClusters []string
	kafkaClusters = append(kafkaClusters, cfg.KAFKAAddr)
	producer := kafka.NewProducer(kafkaClusters, cfg.MailRegistrationTopic)
	authFacade := facade.NewAuthFacade(userService, tokenService, producer)
	authHandler := handler.NewAuthHandler(authFacade)
	mailConsumer := kafka.NewMailConsumer(
		kafka.Config{
			Brokers: []string{cfg.KAFKAAddr},
			GroupID: "my-group",
		},
		cfg.MailRegistrationTopic,
	)
	// --- Router ---
	router := setupRouter(pool, userHandler, authHandler, authMiddleware.Handler())

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
		server:       srv,
		db:           pool,
		mailConsumer: mailConsumer,
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
	go func() {
		slog.Info("starting mail consumer")
		if err := a.mailConsumer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
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
	gin.SetMode(gin.ReleaseMode)
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
