package handlers

import (
	"fmt"
	"net/http"
	"triple-s/storage"
)

// func (b *Bucket) save() error {
// 	filename := "buckets.csv"
// }

func BucketHandler(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	switch r.Method {
	case http.MethodGet:
		GetAllBuckets(w, r, dir)
	case http.MethodPut(w, r, dir):
		PutBucket(w, r, dir, bucket_name)
	case http.MethodDelete:
		DeleteBucket(w, r, dir, bucket_name)
	default:
		panic()
	}
}

func GetAllBuckets(w http.ResponseWriter, r *http.Request, dir string) {
	all_buckets_xml, err := storage.XMLallBuckets(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	fmt.Fprint(w, string(all_buckets_xml))
}

func PutBucket(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	if !valdiate.BucketnameValidation(bucket_name) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
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
}

func DeleteBucket(w http.ResponseWriter, r *http.Request, dir, bucket_name string) {
	err := storage.DeleteBucketStorage(bucket_name, dir)
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
	// w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusNoContent)
}
