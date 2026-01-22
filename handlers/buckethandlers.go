package handlers

import "time"

type Bucket struct {
	Name             string
	CreationTime     time.Time
	LastModifiedTime time.Time

	Content []byte
}

func (b *Bucket) save() error {
	filename := "buckets.csv"

}

func 
