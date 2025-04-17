package config

import (
	"log"
	commonenv "microservice-template/common/env"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SetupDB() (*sqlx.DB, error) {
	dbHost := commonenv.EnvString("DB_HOST", "localhost")
	dbPort := commonenv.EnvString("DB_PORT", "5221")
	dbUser := commonenv.EnvString("DB_USER", "user")
	dbPassword := commonenv.EnvString("DB_PASSWORD", "password")
	dbName := commonenv.EnvString("DB_NAME", "microservice_template_orders_db")

	// connection string
	connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser +
		" password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	// connect to db via sqlx
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}
