package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	storage "triple-s/storage"
	validate "triple-s/validation"
)

// func (b *Bucket) save() error {
// 	filename := "buckets.csv"
// }

func BucketHandler(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	switch r.Method {
	case http.MethodGet:
		GetAllBuckets(w, r, dir)
	case http.MethodPut:
		PutBucket(w, r, dir, bucket_name)
	case http.MethodDelete:
		DeleteBucket(w, r, dir, bucket_name)
	default:
		http.Error(w, "Invalid request(buckets)", http.StatusBadRequest)
	}
}

func GetAllBuckets(w http.ResponseWriter, r *http.Request, dir string) {
	all_buckets_xml, err := storage.XMLallBuckets(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	fmt.Fprint(w, string(all_buckets_xml))
	w.WriteHeader(http.StatusAccepted)
}

func PutBucket(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	if !validate.BucketnameValidation(bucket_name) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}

	err, IsExists := storage.CreateBucket(bucket_name, dir)
	if IsExists {
		http.Error(w, "The requested bucket name is not available."+
			"The bucket namespace is shared by all users of the system. Select a different name and try again.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Date", time.Now().UTC().Format(time.RFC3339))
	w.Header().Set("Location", filepath.Join(dir, bucket_name))
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusCreated)
	// dunno what to add
}

func DeleteBucket(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	err := storage.DeleteBucketStorage(bucket_name, dir)
	if err == nil {
		// w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Date", time.Now().UTC().Format(time.RFC3339))
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err.Error() == "Bucket does not exist" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err.Error() == "Non-empty Bucket" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, "Unconsidered DeleteBucket() case\n Please report to the devs", http.StatusInternalServerError)
}
