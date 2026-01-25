package storage

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"triple-s/structs"
)

func UploadObject(object structs.ObjectMetadata, content io.Reader, bucket_dir string) error {
	object_path := filepath.Join(bucket_dir, object.ObjectKey)
	file, err := os.Create(object_path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	ok, err := IsObjectExist(object, bucket_dir)
	if err != nil && err.Error() != "There's no objects in the bucket" {
		return err
	}

	if !ok {
		return putObjectMetadata(object, bucket_dir)
	} else {
		return EditBucketMetadataTo(object, bucket_dir)
	}
}

func GetObjectContent(object structs.ObjectMetadata, bucket_name, bucket_dir string) (structs.ObjectMetadata, f, error) {
	object_path := filepath.Join(bucket_dir, object.ObjectKey)
	f, err := os.Open(object_path)
	defer f.Close()

	object, err = GetObjectMetadata(object, bucket_dir)
	if err != nil {
		return nil, nil, err
	}
	return object, f, nil
}

func DeleteObjectContent(object structs.ObjectMetadata, bucket_name, bucket_dir string) error {
	objectPath := filepath.Join(bucket_dir, object.ObjectKey)

	if err := os.Remove(objectPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("No such file to delete")
		}
		return err
	}
	return DeleteObjectContent(object, bucket_dir)
}

func putObjectMetadata(object structs.ObjectMetadata, bucket_dir string) error {
	csv_dir := filepath.Join(bucket_dir, "objects.csv")
	f, err := os.OpenFile(csv_dir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err == nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write([]string{
		object.ObjectKey,
		object.ContentLength,
		object.ContentType,
		object.LastModified.Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func IsObjectExist(object structs.ObjectMetadata, bucket_dir string) (bool, error) {
	metadata := filepath.Join("bucket_dir", "objects.csv")
	f, err := os.OpenFile(metadata, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return false, errors.New("There's no objects in the bucket")
	}

	for _, row := range records {
		if row[0] == object.ObjectKey {
			return true, nil
		}
	}

	return false, nil
}

func EditObjectMetadataTo(object_NewMetadata structs.ObjectMetadata, bucket_dir string) error {
	csv_objects := filepath.Join(bucket_dir, "objects.csv")
	f, err := os.OpenFile(csv_objects, os.O_CREATE|os.O_RDONLY, 0644)
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
		if row[0] == object_NewMetadata.ObjectKey {
			row[0] = object_NewMetadata.ObjectKey // optional btw
			row[1] = object_NewMetadata.ContentLength
			row[2] = object_NewMetadata.ContentType
			row[3] = time.Now().Format(time.RFC3339)
			if err := rewriteCSV(csv_objects, records); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("The object's metadata cannot be edit since there's no such")
}

func GetObjectMetadata(object structs.ObjectMetadata, bucket_dir string) (structs.ObjectMetadata, error) {
	csv_objects := filepath.Join(bucket_dir, "objects.csv")
	f, err := os.OpenFile(csv_objects, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)

	records, err := r.ReadAll() // return err if no records
	if err != nil {
		return nil, errors.New("There's no buckets in the storage")
	}

	for _, row := range records {
		// if i == 0 {
		// 	continue
		// }
		if row[0] == object.ObjectKey {
			ContentLength, err := strconv.Atoi(row[1])
			if err != nil {
				return nil, err
			}
			LastModified, err := time.Parse(time.RFC3339, row[3])
			if err != nil {
				return nil, err
			}

			return structs.NewObjectMetadata(row[0], ContentLength, row[2], LastModified), nil
		}
	}
	return nil, errors.New("The object's metadata cannot be edit since there's no such")
}

func DeleteObjectMetadata(object_NewMetadata structs.ObjectMetadata, bucket_dir string) error {
	csv_objects := filepath.Join(bucket_dir, "objects.csv")
	f, err := os.OpenFile(csv_objects, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)

	records, err := r.ReadAll() // return err if no records
	if err != nil {
		return errors.New("There's no buckets in the storage")
	}
	buf := [][]string{}
	for _, row := range records {
		// if i == 0 {
		// 	continue
		// }
		if row[0] != object_NewMetadata.ObjectKey {
			buf = append(buf, row)
		}
	}
	if err := rewriteCSV(csv_objects, records); err != nil {
		return err
	}
	return nil
}
