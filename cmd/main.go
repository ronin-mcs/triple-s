package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(w http.ResponseWriter, r *http.Request, dir string) http.HandlerFunc {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 0 && r.Method == http.MethodGet:
		handlers.BucketHandler(w, r, dir, "")
	case len(parts) == 1 && r.Method != http.MethodGet:
		bucket_name := parts[0]
		handlers.BucketHandler(w, r, dir, bucket_name)
	case len(parts) == 2:
		bucket_name := parts[0]
		object_key := parts[1]
		handlers.ObjectHandler(w, r, dir, bucket_name, object_key)
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

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	handler(w, r, *dataDir)
	// })

	http.HandleFunc("/", makeHandler())

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
