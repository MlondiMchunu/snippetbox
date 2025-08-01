package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	/*Import the models package*/
	"snippetbox.mlodev.net/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Try loading .env file but don't fail if it doesn't exist
	godotenv.Load() // intentionally ignore the errorFatal

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // Default port if not set
	}

	// Initialize logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Log the port used
	infoLog.Printf("Configuring server to listen on port %s", port)

	// Get database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPass := os.Getenv("DB_PASS")
	//caCert := os.Getenv("CA_CERT")

	// Validate all required database configuration
	if dbHost == "" || dbPort == "" || dbPass == "" {
		errorLog.Fatal("Missing database configuration. Please set DB_HOST, DB_PORT, and DB_PASS environment variables")
	}

	// Configure TLS for database connection

	/*tlsConfig := "false" // Default to no TLS
	if caCert != "" {
		rootCertPool := x509.NewCertPool()
		if ok := rootCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			errorLog.Fatal("Failed to parse CA certificate from environment variable")
		}

		// Register custom TLS configuration
		err := mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:    rootCertPool,
			MinVersion: tls.VersionTLS12,
		})
		if err != nil {
			errorLog.Fatal("Failed to register TLS config:", err)
		}
		tlsConfig = "custom"
	}
	*/

	// Database connection setup - using environment variables
	//dsn := "avnadmin:" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/snippetbox?tls=" + tlsConfig + "&parseTime=true"
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//initialize a decoder instance
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Application setup
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConf := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Server configuration - critical for Render compatibility
	srv := &http.Server{
		//for remote db connection
		//Addr: ":" + port, // The colon prefix is required

		//for local db connection
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConf,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server - this log message is important for Render
	/*infoLog.Printf("Starting server on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errorLog.Fatal(err)
	}
	*/

	//local connection
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

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
