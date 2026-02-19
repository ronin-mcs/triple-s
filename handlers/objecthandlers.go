package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	storage "triple-s/storage"
	s "triple-s/structs"
	validate "triple-s/validation"
)

func ObjectHandler(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_key string) {
	object := s.NewObjectMetadata(object_key, r.ContentLength, r.Header.Get("Content-Type"), time.Now())

	switch r.Method {
	case http.MethodPut:
		PutObject(w, r, dir, bucket_name, object)
	case http.MethodGet:
		GetObject(w, r, dir, bucket_name, object)
	case http.MethodDelete:
		DeleteObject(w, r, dir, bucket_name, object)
	default:
		http.Error(w, "Invalid request (object)", http.StatusBadRequest)
	}
}

func PutObject(w http.ResponseWriter, r *http.Request, dir, bucket_name string, object *s.ObjectMetadata) {
	bucket_dir := filepath.Join(dir, bucket_name)
	bucketExistence(w, bucket_name, bucket_dir)

	if !validate.ObejectkeyValidation(object.ObjectKey) {
		http.Error(w, "Invalid object key", http.StatusBadRequest)
	}

	err := storage.UploadObject(object, r.Body, bucket_dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", object.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(object.ContentLength, 10))
	w.WriteHeader(http.StatusOK)
}

func GetObject(w http.ResponseWriter, r *http.Request, dir, bucket_name string, object *s.ObjectMetadata) {
	bucket_dir := filepath.Join(dir, bucket_name)
	bucketExistence(w, bucket_name, bucket_dir)
	objectExistence(w, bucket_dir, object)

	meta, file, err := storage.GetObjectContent(object, bucket_dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(meta.ContentLength, 10))
	w.Header().Set("LastModified", meta.LastModified.Format(time.RFC3339))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}

func DeleteObject(w http.ResponseWriter, r *http.Request, dir, bucket_name string, object *s.ObjectMetadata) {
	bucket_dir := filepath.Join(dir, bucket_name)
	bucketExistence(w, bucket_name, bucket_dir)
	objectExistence(w, bucket_dir, object)

	err := storage.DeleteObjectContent(object, bucket_dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "None")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("LastModified", time.Now().Format(time.RFC3339))
	w.WriteHeader(http.StatusNoContent)
}

func bucketExistence(w http.ResponseWriter, bucket_name, bucket_dir string) {
	ok, err := storage.IsBucketExists(bucket_name, bucket_dir)
	if !ok {
		http.Error(w, "The requested bucket name is not available."+
			"The bucket namespace is shared by all users of the system. Select a different name and try again.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func objectExistence(w http.ResponseWriter, bucket_dir string, object *s.ObjectMetadata) {
	ok, err := storage.IsObjectExist(object, bucket_dir)
	if !ok {
		http.Error(w, "The requested object does not exist", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func extractObjectHeader(r *http.Request) string { // i dont receive Content-Length as int
// 	// since server might want to process it somehow (hypothetically)
// 	return r.Header.Get("Content-Type")
// }
