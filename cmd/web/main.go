package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time" // Added for db.SetConnMaxLifetime

	/*Import the models package*/
	"snippetbox.mlodev.net/internal/models"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Load .env file - will work locally but not break in production
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found - using environment variables")
	}

	// Get port from Render environment variable - THIS IS CRUCIAL FOR RENDER
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // Default port if not specified
	}

	// Initialize logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Access other environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPass := os.Getenv("DB_PASS")
	caCertPath := os.Getenv("CA_CERT_PATH")

	// Validate required configuration
	if dbHost == "" || dbPort == "" || dbPass == "" {
		errorLog.Fatal("Database configuration missing. Set DB_HOST, DB_PORT, and DB_PASS")
	}

	// Load CA certificate if path is specified
	if caCertPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(caCertPath)
		if err != nil {
			errorLog.Fatal("Failed to read CA certificate:", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			errorLog.Fatal("Failed to parse CA certificate")
		}

		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:    rootCertPool,
			MinVersion: tls.VersionTLS12,
		})
		if err != nil {
			errorLog.Fatal(err)
		}
	}

	// Database connection
	dsn := "avnadmin:" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/snippetbox?tls=custom&parseTime=true"
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Setup application
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Configure server - MUST use ":" prefix for Render compatibility
	srv := &http.Server{
		Addr:     ":" + port, // Critical colon prefix
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// IMPORTANT: Log the port binding for Render detection
	infoLog.Printf("SERVER STARTING on port %s", port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute) // Uncommented and fixed

	return db, nil
}
