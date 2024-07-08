package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mreleftheros/snippetbox_ssr/internal/models"
)

type environment struct {
	addr string
}

type application struct {
	infoLog        *log.Logger
	errLog         *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	var env environment
	flag.StringVar(&env.addr, "addr", "localhost:3000", "http address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		errLog.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()
	infoLog.Print("Created connection pool successfully")

	if err = dbPool.Ping(context.Background()); err != nil {
		errLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbPool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	app := &application{
		infoLog:        infoLog,
		errLog:         errLog,
		snippets:       models.NewSnippetModel(dbPool),
		users:          models.NewUserModel(dbPool),
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := http.Server{
		Addr:         env.addr,
		Handler:      app.routes(),
		ErrorLog:     errLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", env.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}
