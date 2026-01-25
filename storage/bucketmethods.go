package storage

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Bucket struct {
	Name             string `xml:"Name"`
	CreationDate     string `xml:"CreationDate"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets Buckets  `xml:"Buckets"`
}

var BucketMap map[string]Bucket

// function to extract file with metadata and error
// function to extract

func XMLallBuckets(dir string) ([]byte, error) {
	f, err := os.OpenFile(dir+"/buckets.csv", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.New("There's no buckets in the storage")
	}

	buckets := []Bucket{}
	for _, row := range records {
		// if i == 0 {
		// 	continue // first row is headers
		// }

		buckets = append(buckets, Bucket{
			Name:             row[0],
			CreationDate:     row[1],
			LastModifiedTime: row[2],
			Status:           row[3],
		})
	}

	result := ListAllMyBucketsResult{
		Buckets: Buckets{Bucket: buckets},
	}

	output, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, errors.New("MarshalIndent() error")
	}
	return output, nil
}

func CreateBucket(bucket_name, dir string) (error, bool) {
	ok, err := IsBucketExists(bucket_name, dir)
	if err != nil {
		return err, false
	}

	bucket_dir := filepath.Join(dir, bucket_name)
	if _, err := os.Stat(bucket_dir); err == nil || !ok {
		return nil, true
	}
	err = os.MkdirAll("data/my-bucket", 0755)
	if err != nil {
		return err, false
	}

	return PutBucketMetadata(bucket_name, time.Now(), time.Now(), "active", dir), false
}

func DeleteBucketStorage(bucket_name, dir string) error {
	ok, err := IsBucketExists(bucket_name, dir)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Bucket does not exist")
	}

	bucket_dir := filepath.Join(dir, bucket_name)
	err = os.Remove(bucket_dir)
	if err != nil {
		return errors.New("Non-empty Bucket")
	}

	return EditBucketMetadataTo(bucket_name, "deleted", dir)
}

func EditBucketMetadataTo(bucket_name string, Status, dir string) error {
	bucket_dir := filepath.Join(dir, "buckets.csv")
	f, err := os.OpenFile(bucket_dir, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)

	records, err := r.ReadAll() // return err if no records
	if err != nil {
		return errors.New("There's no buckets in the storage")
	}

	for _, row := range records {
		// if i == 0 {
		// 	continue
		// }
		if row[0] == bucket_name {
			row[0] = bucket_name // optional btw
			row[2] = time.Now().Format(time.RFC3339)
			row[3] = Status
			if err := rewriteCSV(bucket_dir, records); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("The bucket's metadata cannot be edit since there's no such")
}

func PutBucketMetadata(bucket_name string, CreationTime, LastModifiedTime time.Time, Status, dir string) error {
	bucket_dir := filepath.Join(dir, "buckets.csv")
	f, err := os.OpenFile(bucket_dir,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write([]string{
		bucket_name,
		CreationTime.Format(time.RFC3339),
		LastModifiedTime.Format(time.RFC3339),
		Status,
	})
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func IsBucketExists(bucket_name, dir string) (bool, error) {
	// creating a map and checking by it would be more efficient though
	bucket_dir := filepath.Join(dir, "buckets.csv")
	f, err := os.OpenFile(bucket_dir, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return false, errors.New("There's no buckets in the storage")
	}

	for _, row := range records {
		// if i == 0 {
		// 	continue
		// }
		if row[0] == bucket_name {
			return true, nil
		}
	}
	return false, nil
}

func rewriteCSV(dir string, records [][]string) error {
	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	return w.WriteAll(records)
}
