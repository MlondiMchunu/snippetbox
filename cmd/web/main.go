package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

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
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPass := os.Getenv("DB_PASS")
	caCertPath := os.Getenv("CA_CERT_PATH") // Add this to .env file

	// Get port from Render environment variable, default to 4000
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Initialize logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

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

	// Configure server - using port directly instead of flag for Render compatibility
	srv := &http.Server{
		Addr:     ":" + port, // Use port directly from environment
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
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
	//db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
