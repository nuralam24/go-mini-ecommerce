package database

import (
	"database/sql"
	"os"

	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database/sqlc"
	"go-ecommerce/internal/logger"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *sqlc.Store

func Connect(cfg *config.Config) (*sqlc.Store, error) {
	connStr := getConnectionString()
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if cfg != nil {
		DB.SetMaxOpenConns(cfg.DBMaxOpenConns)
		DB.SetMaxIdleConns(cfg.DBMaxIdleConns)
		DB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
		DB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)
	}

	if err := DB.Ping(); err != nil {
		return nil, err
	}
	Queries = sqlc.NewStore(DB)
	stats := DB.Stats()
	logger.Log.Info().
		Int("db_max_open_conns", stats.MaxOpenConnections).
		Int("db_idle_conns", stats.Idle).
		Msg("Database connected successfully")
	return Queries, nil
}

func Disconnect() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func getConnectionString() string {
	if s := os.Getenv("DATABASE_URL"); s != "" {
		return s
	}
	user := os.Getenv("PGUSER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("PGPASSWORD")
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}
	dbname := os.Getenv("PGDATABASE")
	if dbname == "" {
		dbname = "ecommerce"
	}
	sslmode := os.Getenv("PGSSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
}
