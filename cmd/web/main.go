package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

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
	// Try loading .env file but don't fail if it doesn't exist
	godotenv.Load() // intentionally ignore the error

	// Get port from environment - Render provides this automatically
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // Default port if not set
	}

	// Initialize logging - critical for debugging on Render
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Log the port we're using - helps verify Render configuration
	infoLog.Printf("Configuring server to listen on port %s", port)

	// Get database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPass := os.Getenv("DB_PASS")
	caCertPath := os.Getenv("CA_CERT_PATH")

	// Validate all required database configuration
	if dbHost == "" || dbPort == "" || dbPass == "" {
		errorLog.Fatal("Missing database configuration. Please set DB_HOST, DB_PORT, and DB_PASS environment variables")
	}

	// Load CA certificate if specified - important for secure connections
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
			MinVersion: tls.VersionTLS12, // Enforce modern TLS
		})
		if err != nil {
			errorLog.Fatal(err)
		}
	}

	// Database connection setup - using environment variables
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

	// Application setup
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Server configuration - critical for Render compatibility
	srv := &http.Server{
		Addr:     ":" + port, // The colon prefix is required
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start server - this log message is important for Render
	infoLog.Printf("Starting server on :%s", port)
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

	// Configure connection pool for better performance
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
