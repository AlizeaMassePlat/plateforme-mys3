// storage/storage.go

package storage

import (
	"bytes"
	"context"
	"io"
	"github.com/minio/minio-go/v7" 
)

// nouveau bucket
// func CreateBucket(bucketName string) error {
// 	ctx := context.Background()
// 	err := MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
// 	if err != nil {
// 		exists, errBucketExists := MinioClient.BucketExists(ctx, bucketName)
// 		if errBucketExists == nil && exists {
// 			return errors.New("bucket already exists")
// 		}
// 		return err
// 	}
// 	return nil
// }

// Configuration S3
// const (
// 	s3Endpoint   = "http://minio:9000" 		// endpoint S3
// 	accessKey    = "minioadmin"            	// clé d'accès
// 	secretKey    = "minioadmin"            	// clé secrète
// 	awsRegion    = "us-east-1"             	// Région par défaut
// )

// // Création d'un nouveau bucket
// func CreateBucket(bucketName string) error {
// 	// Préparez l'URL pour la requête PUT
// 	url := fmt.Sprintf("%s/%s", s3Endpoint, bucketName)

// 	// Créez une nouvelle requête HTTP PUT
// 	req, err := http.NewRequest("PUT", url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	// Définir les en-têtes nécessaires
// 	req.Header.Set("Host", strings.Replace(s3Endpoint, "http://", "", 1))
// 	req.Header.Set("x-amz-date", time.Now().UTC().Format("20060102T150405Z"))
// 	req.Header.Set("x-amz-content-sha256", "UNSIGNED-PAYLOAD")

// 	// Ajouter l'authentification
// 	addS3Authorization(req, "PUT", bucketName, "", accessKey, secretKey)

// 	// Créez un client HTTP et envoyez la requête
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Vérifiez le code de statut HTTP
// 	if resp.StatusCode == http.StatusOK {
// 		return nil
// 	} else if resp.StatusCode == http.StatusConflict {
// 		return errors.New("bucket already exists")
// 	}

// 	return fmt.Errorf("failed to create bucket: %s", resp.Status)
// }

// // Ajoute les en-têtes d'autorisation AWS Signature Version 4
// func addS3Authorization(req *http.Request, method, bucket, objectKey, accessKey, secretKey string) {
// 	// Générer des chaînes de date
// 	now := time.Now().UTC()
// 	date := now.Format("20060102")
// 	amzDate := now.Format("20060102T150405Z")

// 	// Créer la chaîne canonique
// 	canonicalRequest := createCanonicalRequest(req, method, bucket, objectKey, amzDate)

// 	// Créer la chaîne de signature
// 	stringToSign := createStringToSign(amzDate, date, awsRegion, canonicalRequest)

// 	// Calculer la signature
// 	signature := calculateSignature(stringToSign, secretKey, date, awsRegion)

// 	// Ajouter l'en-tête d'autorisation
// 	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=%s",
// 		accessKey, date, awsRegion, signature)
// 	req.Header.Set("Authorization", authHeader)
// }

// // Crée la requête canonique
// func createCanonicalRequest(req *http.Request, method, bucket, objectKey, amzDate string) string {
// 	canonicalURI := fmt.Sprintf("/%s", bucket)
// 	canonicalQueryString := ""
// 	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:UNSIGNED-PAYLOAD\nx-amz-date:%s\n", req.Header.Get("Host"), amzDate)
// 	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
// 	payloadHash := "UNSIGNED-PAYLOAD"
// 	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)
// }

// // Crée la chaîne de signature
// func createStringToSign(amzDate, date, region, canonicalRequest string) string {
// 	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", date, region)
// 	hashedCanonicalRequest := hashSHA256([]byte(canonicalRequest))
// 	return fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", amzDate, credentialScope, hashedCanonicalRequest)
// }

// // Calculer la signature
// func calculateSignature(stringToSign, secretKey, date, region string) string {
// 	secret := "AWS4" + secretKey
// 	dateKey := hmacSHA256([]byte(secret), []byte(date))
// 	dateRegionKey := hmacSHA256(dateKey, []byte(region))
// 	dateRegionServiceKey := hmacSHA256(dateRegionKey, []byte("s3"))
// 	signingKey := hmacSHA256(dateRegionServiceKey, []byte("aws4_request"))
// 	signature := hmacSHA256(signingKey, []byte(stringToSign))
// 	return fmt.Sprintf("%x", signature)
// }

// // Fonction HMAC-SHA256
// func hmacSHA256(key, data []byte) []byte {
// 	h := hmac.New(sha256.New, key)
// 	h.Write(data)
// 	return h.Sum(nil)
// }

// // Fonction de hachage SHA256
// func hashSHA256(data []byte) string {
// 	h := sha256.New()
// 	h.Write(data)
// 	return fmt.Sprintf("%x", h.Sum(nil))
// }


func GetBucket(bucketName string) (*Bucket, error) {

	return &Bucket{Name: bucketName}, nil
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
