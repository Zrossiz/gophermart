package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	DBDSN := os.Getenv("DB_DSN")

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
			zap.S().Fatalf("error migrate all files: %v\n", err)
		}
	} else {
		err = migrateOne(db, filenameFlag)
		if err != nil {
			zap.S().Fatalf("error migrate file: %v\n", err)
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
