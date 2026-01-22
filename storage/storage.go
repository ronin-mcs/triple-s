package storage

import (
	"encoding/csv"
	"encoding/xml"
	"os"
)

func XMLallBuckets(dir string) ([]byte, error) {
	f, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.New("There's no buckets in the storage")
	}

	type ListAllMyBucketsResult struct {
		XMLName xml.Name `xml:"ListAllMyBucketsResult"`
		Buckets Buckets  `xml:"Buckets"`
	}

	type Buckets struct {
		Bucket []Bucket `xml:"Bucket"`
	}

	type Bucket struct {
		Name         string `xml:"Name"`
		CreationDate string `xml:"CreationDate"`
	}

	for i, row := range records {
		if i == 0 {
			continue // first row is headers
		}

		buckets = append(buckets, Bucket{
			Name: row[0],
			CreationDate: row[1],
			LastModifiedTime: row[2],
			Status: row[3]
		})

		result := ListAllMyBucketsResult{
			Buckets: Buckets{Bucket: buckets},
		}

		output, err := xml.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, errors.New("MarshalIndent() error")
		}
		return output, nil
	}

}
