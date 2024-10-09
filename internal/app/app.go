package app

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/service"
	"github.com/Zrossiz/gophermart/internal/storage/postgresql"
	"github.com/Zrossiz/gophermart/internal/transport/http/handler"
	"github.com/Zrossiz/gophermart/internal/transport/http/router"
	"github.com/Zrossiz/gophermart/pkg/logger"
	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.Init()
	if err != nil {
		zap.S().Fatalf("config init error: %v", err)
	}

	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("logger int error: %v", err)
	}
	log := zapLogger.ZapLogger

	dbConn, err := postgresql.Connect(cfg.DBDSN)
	if err != nil {
		log.Sugar().Fatalf("erro connect to db: %v", err)
	}
	db := postgresql.New(dbConn)

	s := service.New(service.Storage{
		UserStorage:           &db.UserStore,
		BalanceHistoryStorage: &db.BalanceHistoryStore,
		OrderStorage:          &db.OrderStore,
		TokenStorage:          &db.TokenStore,
		StatusStorage:         &db.StatusStore,
	})

	h := handler.New(handler.Service{
		UserService:           s.UserService,
		BalanceHistoryService: s.BalanceHistoryService,
		OrderService:          s.OrderService,
		StatusService:         s.StatusService,
	})
	r := router.New(h)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
