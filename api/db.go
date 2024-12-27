package api

import (
	"backend/intenal/repository"
	"bytes"
	"database/sql"
	"log"
	"os"
	"unicode/utf8"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
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

	// Hapus BOM jika ada
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})

	// Validasi UTF-8
	if !utf8.Valid(content) {
		// Convert ke UTF-8 jika bukan
		content = bytes.ToValidUTF8(content, []byte(""))
	}

	_, err = app.DB.Connection().Exec(string(content))
	return err
}
