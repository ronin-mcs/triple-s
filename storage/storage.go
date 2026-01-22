package storage

import "os"

func XMLallBuckets(dir string) error {
	f, err := os.Open(dir + "/buckets.csv")
	if err != nil {
		return err
	}
	defer f.Close()

}
