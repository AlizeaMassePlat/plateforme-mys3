package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func AjoutBucket(bucketName string) error {
	// Define the endpoint and access keys
	endpoint := "http://localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	region := "us-east-1"
	service := "s3"

	// Build the request URL
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	// Create HTTP PUT request to create the bucket
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	// Add headers for authentication and S3 compatibility
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	date := time.Now().UTC().Format("20060102")
	req.Header.Add("x-amz-date", amzDate)
	req.Header.Add("Host", "localhost:9000")
	req.Header.Add("Content-Length", "0")

	// Generate the signature
	signature := generateSignature(req, accessKeyID, secretAccessKey, date, region, service)
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/%s/aws4_request, SignedHeaders=host;x-amz-date, Signature=%s",
		accessKeyID, date, region, service, signature)
	req.Header.Add("Authorization", authHeader)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create bucket: %s", resp.Status)
	}

	return nil
}

func generateSignature(req *http.Request, accessKeyID, secretAccessKey, date, region, service string) string {
	canonicalURI := req.URL.Path
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-date:%s\n", req.Host, req.Header.Get("x-amz-date"))
	signedHeaders := "host;x-amz-date"
	payloadHash := sha256.Sum256([]byte{}) // Hash for an empty payload

	// Create canonical request
	canonicalRequest := fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s", req.Method, canonicalURI, "", canonicalHeaders, signedHeaders, hex.EncodeToString(payloadHash[:]))

	// Create string to sign
	canonicalRequestHash := sha256.Sum256([]byte(canonicalRequest))
	canonicalRequestHashStr := hex.EncodeToString(canonicalRequestHash[:])
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/%s/%s/aws4_request\n%s", req.Header.Get("x-amz-date"), date, region, service, canonicalRequestHashStr)

	// Calculate the signature
	kDate := hmacSHA256([]byte("AWS4"+secretAccessKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	signature := hmacSHA256(kSigning, stringToSign)

	return hex.EncodeToString(signature)
}


// hmacSHA256 calculates an HMAC-SHA256 hash
func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}
