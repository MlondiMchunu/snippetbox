package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
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
	caCertPath := os.Getenv("CA_CERT_PATH") // Add this to your .env file

	//fmt.Printf("Database host: %s, port: %s\n", dbHost, dbPort)

	//define new cmd line flag for addr

	//local connection
	//addr := flag.String("addr", ":4000", "HTTP network address")

	//managed db connection
	addr := flag.String("addr", ":4000", "HTTP network address")

	//local db connection
	//dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	//managed db connection
	dsn := flag.String("dsn", "avnadmin:"+dbPass+"@tcp("+dbHost+":"+dbPort+")/snippetbox?tls=custom&parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Load CA certificate
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caCertPath)
	if err != nil {
		errorLog.Fatal("Failed to read CA certificate:", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		errorLog.Fatal("Failed to parse CA certificate")
	}

	// Register custom TLS config
	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:    rootCertPool,
		MinVersion: tls.VersionTLS12, // Enforce minimum TLS 1.2
	})
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	//initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//add templateCache to the application dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db}, //initialize a models.SnippetModel instance & add it to app dependencies
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	/*register routes without declaring a servemux
	*NB avoid on production apps for security reasons
	 */
	//http.HandleFunc("/", home)

	//port := 4000

	infoLog.Printf("Starting server on %s ", *addr)

	err = srv.ListenAndServe()

	/*part of registering routes withut declaring a servemux*/
	//err := http.ListenAndServe(":4000,nil")
	// errorLog.Fatal(err)

	//find . -name "*.go" | entr -r sh -c 'echo "== Restarting =="; go run ./cmd/web'
}

// openDB function wraps sql.Open() & returns a sql.DB
// connetion pool for a given  dsn
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
