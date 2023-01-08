package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VtG242/boss/internal/models"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// Define app struct to keep global settings
type application struct {
	log       *log.Logger
	players   *models.PlayersModel
	templates map[string]*template.Template
}

func main() {
	// logrus global settings -  if not set some messages logged in default format
	log.SetFormatter(&log.JSONFormatter{})

	// define command-line flags
	addr := flag.String("addr", ":8300", "HTTP network address")
	dsn := flag.String("dsn", "boss:password@/BOSS?parseTime=true", "MySQL data source name")
	flag.Parse()

	// database pool init
	db_pool, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db_pool.Close()

	// Initialize template cache
	templates_cache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// initialize global app dependencies
	app := &application{
		log: &log.Logger{
			Formatter: &log.JSONFormatter{},
			Level:     log.DebugLevel,
			Out:       os.Stderr,
			//ReportCaller: true,
		},
		players:   &models.PlayersModel{Pool: db_pool},
		templates: templates_cache,
	}

	// Initialize a new http.Server struct
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	// Run server in a goroutine so that it doesn't block.
	log.Info("BOSS Admin server - Start listening at port ", *addr)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// e.g. port used
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C), SIGTERM,
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	log.Warn("BOSS Admin server: Attempt to perform gracefull shutdown")
	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("BOSS Admin server: graceful shutdown performed succesufly")
	os.Exit(0)

}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
