package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/gophermart/internal/api"
	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/service"
	"github.com/Zrossiz/gophermart/internal/storage/postgresql"
	"github.com/Zrossiz/gophermart/internal/transport/http/handler"
	"github.com/Zrossiz/gophermart/internal/transport/http/router"
	"github.com/Zrossiz/gophermart/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.Init()
	if err != nil {
		fmt.Println(err)
		zap.S().Fatalf("config init error: %v", err)
	}

	zapLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err)
		zap.S().Fatalf("logger int error: %v", err)
	}
	log := zapLogger.ZapLogger

	a := api.New(cfg)

	dbConn, err := postgresql.Connect(cfg.DBDSN)
	if err != nil {
		log.Sugar().Fatalf("error connect to db: %v", err)
	}
	db := postgresql.New(dbConn, log)

	s := service.New(service.Storage{
		UserStorage:           &db.UserStore,
		BalanceHistoryStorage: &db.BalanceHistoryStore,
		OrderStorage:          &db.OrderStore,
		TokenStorage:          &db.TokenStore,
		StatusStorage:         &db.StatusStore,
	}, cfg, log, a)

	h := handler.New(handler.Service{
		UserService:           s.UserService,
		BalanceHistoryService: s.BalanceHistoryService,
		OrderService:          s.OrderService,
		StatusService:         s.StatusService,
	})
	r := router.New(h)

	cr := cron.New()
	cr.AddFunc("*/1 * * *", func() {
		log.Info("starting cron...")
		s.OrderService.UpdateOrders()
		log.Info("cron ended")
	})

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	log.Sugar().Infof("Starting server on addr: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
