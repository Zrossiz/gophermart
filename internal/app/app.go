package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Zrossiz/gophermart/internal/api"
	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/service"
	"github.com/Zrossiz/gophermart/internal/storage/postgresql"
	"github.com/Zrossiz/gophermart/internal/transport/http/handler"
	"github.com/Zrossiz/gophermart/internal/transport/http/router"
	"github.com/Zrossiz/gophermart/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func Start() {
	cfg, err := config.Init()
	if err != nil {
		fmt.Println(err)
		zap.S().Fatalf("config init error: %v", err)
	}

	if len(cfg.AutoMigrate) == 0 {
		migrate()
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

// TODO: delete auto migrate after complete project
func migrate() {
	_ = godotenv.Load()
	DBDSN := os.Getenv("DATABASE_URI")
	fmt.Printf("start migrate: %s\n", DBDSN)
	db, err := pgxpool.Connect(context.Background(), DBDSN)
	if err != nil {
		fmt.Printf("error connecting to db: %v\n", err)
		return
	}
	defer db.Close()

	var filenameFlag string
	flag.StringVar(&filenameFlag, "f", "", "filename for migrated file")
	flag.Parse()

	if filenameFlag == "" {
		err = migrateAll(db)
		if err != nil {
			fmt.Printf("error migrate all files: %v\n", err)
		}
	} else {
		err = migrateOne(db, filenameFlag)
		if err != nil {
			fmt.Printf("error migrate file: %v\n", err)
		}
	}

	fmt.Println("Schema created successfully!")
}

func migrateAll(db *pgxpool.Pool) error {
	files, err := os.ReadDir("migration")
	if err != nil {
		return err
	}

	for _, f := range files {
		sqlFilePath := filepath.Join("migration", f.Name())
		schemaSQL, err := os.ReadFile(sqlFilePath)
		if err != nil {
			return err
		}

		_, err = db.Exec(context.Background(), string(schemaSQL))
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateOne(db *pgxpool.Pool, filename string) error {
	sqlFilePath := filepath.Join("migration", filename)
	schemaSQL, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		return err
	}

	return nil
}
