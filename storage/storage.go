// storage/storage.go

package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"github.com/minio/minio-go/v7" 
)

// nouveau bucket
func CreateBucket(bucketName string) error {
	ctx := context.Background()
	err := MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return errors.New("bucket already exists")
		}
		return err
	}
	return nil
}


func GetBucket(bucketName string) (*Bucket, error) {

	return &Bucket{Name: bucketName}, nil
}

// Liste tous les buckets

func ListBuckets() ([]Bucket, error) {
	ctx := context.Background()
	bucketsList, err := MinioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var result []Bucket
	for _, b := range bucketsList {
		bucket := Bucket{Name: b.Name}
		result = append(result, bucket)
	}
	return result, nil
}

// Ajoute un objet dans un bucket

func (b *Bucket) PutObject(objectName string, data []byte) error {
	ctx := context.Background()
	dataReader := bytes.NewReader(data)

	_, err := MinioClient.PutObject(ctx, b.Name, objectName, dataReader, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Récupère un objet dans un bucket

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	ctx := context.Background()

	object, err := MinioClient.GetObject(ctx, b.Name, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Supprime un bucket

func DeleteBucket(bucketName string) error {
	ctx := context.Background()

	err := MinioClient.RemoveBucket(ctx, bucketName)
	if err != nil {
		return err
	}

	return nil
}

// Supprime un objet dans un bucket

func (b *Bucket) DeleteObject(objectName string) error {
	ctx := context.Background()

	err := MinioClient.RemoveObject(ctx, b.Name, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
