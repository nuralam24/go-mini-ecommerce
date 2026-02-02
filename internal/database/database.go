package database

import (
	"database/sql"
	"log"
	"os"

	"go-ecommerce/internal/database/sqlc"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var Queries *sqlc.Store

func Connect() error {
	// DATABASE_URL is loaded from env (e.g. postgres://user:pass@localhost:5432/dbname?sslmode=disable)
	connStr := getConnectionString()
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err := DB.Ping(); err != nil {
		return err
	}
	Queries = sqlc.NewStore(DB)
	log.Println("Database connected successfully")
	return nil
}

func Disconnect() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func getConnectionString() string {
	// Prefer DATABASE_URL
	if s := os.Getenv("DATABASE_URL"); s != "" {
		return s
	}
	// Fallback: build from individual vars
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
