package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utint"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func (cfg *Config) InitPostgres() serror.SError {
	sqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname =%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", sqlConn)
	if err != nil {
		log.Fatalf("Failed connect to database %+v", err)
		return serror.NewFromError(err)
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(utint.StringToInt(utstring.Env("DB_CONNECTION_LIFETIME", "15"), 15)))
	db.SetMaxIdleConns(int(utint.StringToInt(utstring.Env("DB_CONN_MAX_IDLE", "5"), 5)))
	db.SetMaxOpenConns(int(utint.StringToInt(utstring.Env("DB_CONN_MAX_OPEN", "0"), 0)))

	cfg.DB = db

	GlobalShutdown.RegisterGracefullyShutdown("database/postgres", func(ctx context.Context) error {
		return cfg.DB.Close()
	})

	return nil
}
