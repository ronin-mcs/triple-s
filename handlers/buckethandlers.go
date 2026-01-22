package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type Bucket struct {
	Name             string
	CreationTime     time.Time
	LastModifiedTime time.Time

	Content []byte
}

func (b *Bucket) save() error {
	filename := "buckets.csv"

}

func GetAllBuckets(w http.ResponseWriter, r *http.Request, dir string) {
	all_buckets_xml, err := storage.XMLallBuckets(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	fmt.Fprint(w, string(all_buckets_xml))
}

func PutBucket(w http.ResponseWriter, r *http.Request, dir string, bucket_name string) {
	
}
