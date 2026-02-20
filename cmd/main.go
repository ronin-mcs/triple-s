package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"triple-s/handlers"
)

var validPath = regexp.MustCompile(`^/([a-zA-Z0-9]+)(/[a-zA-Z0-9]+)*$`)
func Handler(w http.ResponseWriter, r *http.Request, dir string) {

	if _, err := url.Parse(r.URL.Path); err != nil {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1 && parts[0] == "" && r.Method == http.MethodGet:
		handlers.BucketHandler(w, r, dir, "")
	case len(parts) == 1 && r.Method != http.MethodGet:
		bucket_name := parts[0]
		handlers.BucketHandler(w, r, dir, bucket_name)
	case len(parts) == 2:
		bucket_name := parts[0]
		object_key := parts[1]
		handlers.ObjectHandler(w, r, dir, bucket_name, object_key)
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("1. Handler started: %s", r.URL.Path)

	bucket_name := strings.TrimPrefix(r.URL.Path, "/buckets/")
	log.Printf("2. Bucket name: %q", bucket_name)

	bucket_dir := filepath.Join("./data", bucket_name)
	log.Printf("3. Bucket dir: %q", bucket_dir)

	err := os.MkdirAll(bucket_dir, 0755)
	log.Printf("4. MkdirAll err: %v", err)

	w.WriteHeader(201) // Должно дойти сюда
	log.Println("5. Response sent")
}

func main() {
	flag.Usage = func() {
		fmt.Println("Simple Storage Service.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  triple-s [-port N] [-dir S]")
		fmt.Println("  triple-s --help")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	// ===================================================================
	port := flag.Int("port", 8080, "Port number")
	dir := flag.String("dir", "./data", "Path to directory")

	if *port <= 0 || *port > 65535 {
		log.Fatalf("invalid port: %d", *port)
	}

	if *dir == "" {
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(*dir)
	if os.IsNotExist(err) {
		log.Fatalf("directory does not exist: %s", *dir)
	}
	if err == nil && !info.IsDir() {
		log.Fatalf("path is not a directory: %s", *dir)
	}

	// ===================================================================

	flag.Parse()

	// http.HandleFunc("/", handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, *dir)
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
