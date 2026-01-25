package handlers

func PutObject(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_key string) {
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

	ContentType := extractObjectHeader(r)
	object := structs.NewObjectMetadata(object_key, r.ContentLength, ContentType, time.Now())
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

func GetObject(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_key string) {
	err := 
}

func extractObjectHeader(r *http.Request) string { // i dont receive Content-Length as int
	// since server might want to process it somehow (hypothetically)
	return r.Header.Get("Content-Type")
}
