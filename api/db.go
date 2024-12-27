package api

import (
	"backend/intenal/repository"
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"strings"
)

type Application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	Auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) ConnectToDB() (*sql.DB, error) {
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")
	return connection, nil
}

func (app *Application) RunMigrations() error {
	content, err := os.ReadFile("movies.sql")
	if err != nil {
		return err
	}

	// Pecah menjadi pernyataan SQL individual
	statements := strings.Split(string(content), ";")

	// Eksekusi setiap pernyataan
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err = app.DB.Connection().Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
