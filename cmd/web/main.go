package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	//define new cmd line flag for addr
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	//use the http.NewServeMux() i.e router function to initialize a new servemux
	mux := http.NewServeMux()

	//create a file server which serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//use the mux.Handle() function to register the file server as the handler for
	//all URL paths that start with "/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//manually register handler function with the servemux
	mux.HandleFunc("/", http.HandlerFunc(app.home))

	//transforms function to a handler and registers it in one step
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	/*register routes without declaring a servemux
	*NB avoid on production apps for security reasons
	 */
	//http.HandleFunc("/", home)

	//port := 4000

	infoLog.Printf("Starting server on %s ", *addr)

	err := srv.ListenAndServe()

	/*part of registering routes withut declaring a servemux*/
	//err := http.ListenAndServe(":4000,nil")
	errorLog.Fatal(err)

	//find . -name "*.go" | entr -r sh -c 'echo "== Restarting =="; go run ./cmd/web'
}
