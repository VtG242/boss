package main

import (
	"database/sql"
	"flag"
	"html/template"
	"net/http"
	"os"

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
		players:        &models.PlayersModel{Pool: db_pool},
		templates: templates_cache,
	}

	// Initialize a new http.Server struct
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Info("Starting server on ", *addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
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
