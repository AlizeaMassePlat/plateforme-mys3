package storage

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/minio/minio-go/v7"
)

// Représentation XML pour la réponse de liste de buckets
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Xmlns   string   `xml:"xmlns,attr"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

// Représente un bucket
type Bucket struct {
	XMLName      xml.Name `xml:"Bucket"`
	Name         string   `xml:"Name"`
	CreationDate string   `xml:"CreationDate"` // Assurez-vous que CreationDate est défini
}

// Crée un bucket
func CreateBucket(bucketName string) error {
	ctx := context.Background()
	err := MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		}
		return err
	}
	return nil
}

// Vérifie si un bucket existe déjà
func BucketExists(bucketName string) (bool, error) {
	ctx := context.Background()
	return MinioClient.BucketExists(ctx, bucketName)
}

// Liste tous les buckets
func ListBuckets() (*ListAllMyBucketsResult, error) {
	ctx := context.Background()
	bucketsList, err := MinioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	result := &ListAllMyBucketsResult{
		Xmlns: "http://s3.amazonaws.com/doc/2006-03-01/",
	}

	for _, b := range bucketsList {
		bucket := Bucket{
			Name:         b.Name,
			CreationDate: b.CreationDate.Format(time.RFC3339), // Utilise le champ CreationDate
		}
		result.Buckets = append(result.Buckets, bucket)
	}

	return result, nil
}

// Convertit ListAllMyBucketsResult en XML
func (b *ListAllMyBucketsResult) ToXML() ([]byte, error) {
	return xml.MarshalIndent(b, "", "  ")
}
