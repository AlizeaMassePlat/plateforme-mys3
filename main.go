package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"my-s3-clone/storage"
)

func main() {
	storage.InitMinioClient("minio:9000", "minioadmin", "minioadmin")

	http.HandleFunc("/buckets", handleBuckets)
	http.HandleFunc("/buckets/", handleObjects)
	http.ListenAndServe(":8080", nil)
}

func handleBuckets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var bucket storage.Bucket
		body, _ := io.ReadAll(r.Body) 
		xml.Unmarshal(body, &bucket)

		err := storage.CreateBucket(bucket.Name) 
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated) 

	case http.MethodGet:
		buckets, err := storage.ListBuckets() 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(buckets)

	case http.MethodDelete:
		bucketName := r.URL.Path[len("/buckets/"):]
		err := storage.DeleteBucket(bucketName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}


func handleObjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		bucketName := r.URL.Path[len("/buckets/"):]
		objectName := r.URL.Query().Get("object")

		bucket, err := storage.GetBucket(bucketName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data, _ := io.ReadAll(r.Body) 
		err = bucket.PutObject(objectName, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) 
	case http.MethodGet:
		bucketName := r.URL.Path[len("/buckets/"):]
		objectName := r.URL.Query().Get("object")

		bucket, err := storage.GetBucket(bucketName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		object, err := bucket.GetObject(objectName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(object)

	case http.MethodDelete:
		bucketName := r.URL.Path[len("/buckets/"):]
		objectName := r.URL.Query().Get("object")

		bucket, err := storage.GetBucket(bucketName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = bucket.DeleteObject(objectName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent) 
	}
}
