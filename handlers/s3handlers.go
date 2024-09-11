package handlers

import (
    "io"
    "my-s3-clone/storage"
    "my-s3-clone/dto"
    "net/http"
    "github.com/gorilla/mux"
    "log"
    "time"
    "encoding/xml"
    "fmt"
    "os"
    "strconv"
)

// List all buckets
func HandleListBuckets(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        log.Println("Calling storage.ListBuckets to get the list of buckets.")
        buckets := s.ListBuckets()

        if len(buckets) == 0 {
            log.Println("No buckets found in storage.")
        } else {
            log.Printf("Found %d buckets.", len(buckets))
        }

        var bucketList []dto.Bucket
        for _, bucketName := range buckets {
            log.Printf("Adding bucket: %s", bucketName)
            bucketList = append(bucketList, dto.Bucket{
                Name:         bucketName,
                CreationDate: time.Now(),
            })
        }

        response := dto.ListAllMyBucketsResult{
            Buckets: bucketList,
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)

        log.Println("Encoding response as XML and sending it.")
        if err := xml.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Erreur lors de l'encodage des buckets en XML", http.StatusInternalServerError)
            log.Printf("Erreur lors de l'encodage des buckets: %v", err)
        }
    }
}

// Create a bucket
func HandleCreateBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        if r.Method != "PUT" {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        err := s.CreateBucket(bucketName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        bucketResponse := dto.ListAllMyBucketsResult{
            Buckets: []dto.Bucket{
                {
                    Name:         bucketName,
                    CreationDate: time.Now(),
                },
            },
        }

        w.Header().Set("Content-Type", "application/xml")
        w.Header().Set("location", r.URL.String())
        w.WriteHeader(http.StatusOK)
        if err := xml.NewEncoder(w).Encode(bucketResponse); err != nil {
            http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
        }
    }
}

// Get bucket info or location
func HandleGetBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        log.Printf("Requête GET pour le bucket: %s", bucketName)

        locationParam := r.URL.Query().Get("location")
        log.Printf("Location Param: %s", locationParam)

        if locationParam != "" {
            log.Printf("Demande de localisation pour le bucket: %s", bucketName)
            w.Header().Set("Content-Type", "application/xml")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
            return
        }

        exists, err := s.CheckBucketExists(bucketName)
        if err != nil || !exists {
            if err != nil {
                log.Printf("Erreur lors de la vérification du bucket: %v", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            } else {
                log.Printf("Bucket non trouvé: %s", bucketName)
                http.Error(w, "Bucket not found", http.StatusNotFound)
            }
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Bucket exists and is accessible."))
    }
}

// Add an object
func HandleAddObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        if bucketName == "" || objectName == "" {
            http.Error(w, "Bucket name and object name are required", http.StatusBadRequest)
            return
        }

        err := s.AddObject(bucketName, objectName, r.Body, r.Header.Get("X-Amz-Content-Sha256"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Upload successful"))
        log.Printf("Successfully uploaded object: %s in bucket: %s", objectName, bucketName)
    }
}

// Check if an object exists
func HandleCheckObjectExist(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        if bucketName == "" || objectName == "" {
            http.Error(w, "Bucket name and object name are required", http.StatusBadRequest)
            return
        }

        exists, lastModified, size, err := s.CheckObjectExist(bucketName, objectName)
        if err != nil || !exists {
            if !exists {
                http.Error(w, "Object not found", http.StatusNotFound)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
        w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
        w.WriteHeader(http.StatusOK)
    }
}

// Download an object
func HandleDownloadObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        data, fileInfo, err := s.GetObject(bucketName, objectName)
        if err != nil {
            if os.IsNotExist(err) {
                http.Error(w, "File not found", http.StatusNotFound)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/octet-stream")
        w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
        w.Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))
        w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
        w.WriteHeader(http.StatusOK)

        if _, err := w.Write(data); err != nil {
            http.Error(w, "Failed to write file content", http.StatusInternalServerError)
        }
    }
}

// List objects in a bucket
func HandleListObjects(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        queryParams := r.URL.Query()
        prefix := queryParams.Get("prefix")
        marker := queryParams.Get("marker")
        maxKeys := queryParams.Get("max-keys")

        if maxKeys == "" {
            maxKeys = "1000"
        }

        maxKeysInt, err := strconv.Atoi(maxKeys)
        if err != nil {
            http.Error(w, "Invalid max-keys value", http.StatusBadRequest)
            return
        }

        objects, err := s.ListObjects(bucketName, prefix, marker, maxKeysInt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        if err := xml.NewEncoder(w).Encode(objects); err != nil {
            http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
        }
    }
}

// Delete a bucket
func HandleDeleteBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        if bucketName == "" {
            http.Error(w, "Bucket name is required", http.StatusBadRequest)
            return
        }

        err := s.DeleteBucket(bucketName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

// Delete a single object
func HandleDeleteObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received DELETE request: %s %s", r.Method, r.URL.Path)
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        err := s.DeleteObject(bucketName, objectName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        deleteResult := dto.DeleteResult{
            DeletedResult: []dto.Deleted{
                {Key: objectName},
            },
        }

        response, err := xml.Marshal(deleteResult)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
    }
}

// Batch delete objects
func HandleDeleteBatch(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received POST ?delete request for batch deletion: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading request body", http.StatusInternalServerError)
            return
        }

        var deleteReq dto.DeleteBatchRequest
        err = xml.Unmarshal(body, &deleteReq)
        if err != nil {
            http.Error(w, "Error parsing XML", http.StatusInternalServerError)
            return
        }

        var deletedObjects []dto.Deleted
        for _, objectToDelete := range deleteReq.Objects {
            err := s.DeleteObject(bucketName, objectToDelete.Key)
            if err != nil {
                http.Error(w, "Error deleting object", http.StatusInternalServerError)
                return
            }

            deletedObjects = append(deletedObjects, dto.Deleted{Key: objectToDelete.Key})
        }

        deleteResult := dto.DeleteResult{
            DeletedResult: deletedObjects,
        }

        response, err := xml.Marshal(deleteResult)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
    }
}
