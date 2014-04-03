package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	port string
	path string
)

func main() {
	// Catch flags
	flag.StringVar(&port, "port", "8000", "Port to serve content on")
	flag.StringVar(&path, "path", ".", "Directory to serve content from")
	flag.Parse()

	// We want the hostname to make a nice easy to copy link
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "0.0.0.0"
	}
	hostname = strings.ToLower(hostname)

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		absolutePath = "./"
	}

	log.Printf("Starting http server to serve '%s' at:\nhttp://%s:%s/", absolutePath, hostname, port)
	fileHandler := http.FileServer(http.Dir(absolutePath))
	wrappedHandler := AccessLoggingHandler(fileHandler)
	log.Fatal(http.ListenAndServe(":"+port, wrappedHandler))
}

// A bit excessive? Perhaps. But I enjoy seeing what's going on while the server is running.
func AccessLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RemoteAddr, "=>", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
