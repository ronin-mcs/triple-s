package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 0 && r.Method == http.MethodGet:
		bucket_name := parts[0]
		storage.XMLallBuckets()
	case 1:

	case 2:

	default:
		// какая-то ошибка, хз пока какая

	}

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	port := flag.String("port", "8080", "Port number")
	dir := flag.String("dir", "./data", "Path to directory")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, *dataDir)
	})

	http.HandleFunc("/", makeHandler())

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
