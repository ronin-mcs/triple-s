package structs

import "time"

type ObjectMetadata struct {
	ObjectKey     string    `xml:"ObjectKey"`
	ContentLength int64     `xml:"ContentLength"`
	ContentType   string    `xml:"ContentType"`
	LastModified  time.Time `xml:"LastModified"`
}

func NewObjectMetadata(ObjectKey string, ContentLength int64, ContentType string, LastModified time.Time) *ObjectMetadata {
	return &ObjectMetadata{
		ObjectKey:     ObjectKey,
		ContentLength: ContentLength,
		ContentType:   ContentType,
		LastModified:  LastModified,
	}
}
