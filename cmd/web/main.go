package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mreleftheros/snippetbox_ssr/internal/models"
	"log"
	"net/http"
	"os"
)

type environment struct {
	addr string
}

type application struct {
	infoLog  *log.Logger
	errLog   *log.Logger
	snippets *models.SnippetModel
}

func main() {
	var env environment
	flag.StringVar(&env.addr, "addr", "localhost:5000", "http address")
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

	app := &application{
		infoLog:  infoLog,
		errLog:   errLog,
		snippets: models.NewSnippetModel(dbPool),
	}

	srv := http.Server{
		Addr:     env.addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting server on %s", env.addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
