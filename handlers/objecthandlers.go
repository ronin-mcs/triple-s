package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func ObjectHandler(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_key string) {
	ContentType := extractObjectHeader(r)
	object := structs.NewObjectMetadata(object_key, r.ContentLength, ContentType, time.Now())

	switch r.Method {
	case http.MethodPut:
		PutObject(w, r, dir, bucket_name, object)
	case http.MethodGet:
		GetObject(w, r, dir, bucket_name, object)
	}
}

func PutObject(w http.ResponseWriter, r *http.Request, dir, bucket_name string, object structs.ObjectMetadata) {
	ok, err := !storage.IsBucketExists(bucket_name, dir)
	if !ok {
		http.Error(w, "The requested bucket name is not available."+
			"The bucket namespace is shared by all users of the system. Select a different name and try again.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !validate.ObejectkeyValidation(object_key) {
		http.Error(w, "Invalid object key", http.StatusBadRequest)
	}

	bucketdir := filepath.Join(dir, bucket_name)

	err = storage.UploadObject(object, r.Body, bucketdir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", object.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(object.ContentLength, 10))
	w.WriteHeader(http.StatusOK)
}

func GetObject(w http.ResponseWriter, r *http.Request, dir, bucket_name, object stringstructs.ObjectMetadata) {
	ok, err := !storage.IsBucketExists(bucket_name, dir)
	if !ok {
		http.Error(w, "The requested bucket name is not available."+
			"The bucket namespace is shared by all users of the system. Select a different name and try again.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bucket_dir := filepath.Join(dir, bucket_name)
	ok, err = storage.IsObjectExist(object, bucket_dir)
	if !ok {
		http.Error(w, "The requested object does not exist", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	meta, file, err := storage.GetObjectContent(object, bucket_name, bucket_dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(meta.Size, 10))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}

func extractObjectHeader(r *http.Request) string { // i dont receive Content-Length as int
	// since server might want to process it somehow (hypothetically)
	return r.Header.Get("Content-Type")
}
