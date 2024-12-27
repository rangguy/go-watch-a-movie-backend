package main

import (
	"backend/api"
	"backend/intenal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

func main() {
	// set Application config
	var app api.Application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=postgres.railway.internal port=5432 user=postgres password=CwBXIKpmrDqkNjJpAAOEKOxFGxKIdKnn dbname=railway sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	//flag.StringVar(&app.DSN, "dsn", "host=localhost port=5433 user=postgres password=postgresql dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "84ab2084877e4a8514a278ff851c9976", "api key")
	flag.Parse()

	// connect to database
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	// Jalankan migrasi setelah koneksi DB
	if err := app.RunMigrations(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	app.Auth = api.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting Application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
