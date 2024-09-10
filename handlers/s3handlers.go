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

// Liste des buckets
func HandleListBuckets(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    log.Println("Calling storage.ListBuckets to get the list of buckets.")
    buckets := storage.ListBuckets()

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


// création d'un bucket
func HandleCreateBucket(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    if r.Method != "PUT" {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    vars := mux.Vars(r)
    bucketName := vars["bucketName"]

    err := storage.CreateBucket(bucketName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // log.Printf("Route appelée pour bucket: %s, méthode: %s", bucketName, r.Method)

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
        // log.Printf("Erreur lors de l'encodage XML: %v", err)
        http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
    }
}



// Gestion des requêtes GET sur un bucket
func HandleGetBucket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bucketName := vars["bucketName"]

    log.Printf("Requête GET pour le bucket: %s", bucketName)

    // Get the location query parameter
    locationParam := r.URL.Query().Get("location")
    log.Printf("Location Param: %s", locationParam)

    // Only process the location request if the location parameter is not empty
    if locationParam != "" {
        log.Printf("Demande de localisation pour le bucket: %s", bucketName)

        // Respond with location constraint in XML format
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
        return
    } else {
        log.Printf("Pas de demande de localisation pour le bucket: %s", bucketName)
    }

    if exists, err := storage.CheckBucketExists(bucketName); err != nil || !exists {
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


func HandleAddObject(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    vars := mux.Vars(r)
    bucketName := vars["bucketName"]
    objectName := vars["objectName"]

    if bucketName == "" {
        http.Error(w, "Bucket name is required", http.StatusBadRequest)
        log.Println("Bucket name is missing")
        return
    }

    if objectName == "" {
        http.Error(w, "Object name is required", http.StatusBadRequest)
        log.Println("Object name is missing")
        return
    }

    err := storage.AddObject(bucketName, objectName, r.Body, r.Header.Get("X-Amz-Content-Sha256"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error while uploading object: %v", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Upload successful"))
    log.Printf("Successfully uploaded object: %s in bucket: %s", objectName, bucketName)
}

func HandleCheckObjectExist(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    vars := mux.Vars(r)
    bucketName := vars["bucketName"]
    objectName := vars["objectName"]

    if bucketName == "" {
        http.Error(w, "Bucket name is required", http.StatusBadRequest)
        log.Println("Bucket name is missing")
        return
    }

    if objectName == "" {
        http.Error(w, "Object name is required", http.StatusBadRequest)
        log.Println("Object name is missing")
        return
    }

    exists, lastModified, size, err := storage.CheckObjectExist(bucketName, objectName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error while checking object existence: %v", err)
        return
    }

    if !exists {
        http.Error(w, "Object not found", http.StatusNotFound)
        log.Printf("Object not found: %s", objectName)
        return
    }

    w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
    w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
    w.WriteHeader(http.StatusOK)
    log.Printf("Object %s exists in bucket %s", objectName, bucketName)
}

func HandleDownloadObject(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)

    vars := mux.Vars(r)
    bucketName := vars["bucketName"]
    objectName := vars["objectName"]

    if bucketName == "" {
        http.Error(w, "Bucket name is required", http.StatusBadRequest)
        log.Println("Bucket name is missing")
        return
    }

    if objectName == "" {
        http.Error(w, "Object name is required", http.StatusBadRequest)
        log.Println("Object name is missing")
        return
    }

    data, fileInfo, err := storage.GetObject(bucketName, objectName)
    if err != nil {
        if os.IsNotExist(err) {
            http.Error(w, "File not found", http.StatusNotFound)
            log.Printf("Object %s in bucket %s not found", objectName, bucketName)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            log.Printf("Error while retrieving object: %v", err)
        }
        return
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
    w.Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))
    w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    w.WriteHeader(http.StatusOK)

    // Écriture du contenu du fichier dans la réponse
    if _, err := w.Write(data); err != nil {
        http.Error(w, "Failed to write file content", http.StatusInternalServerError)
        log.Printf("Failed to write file content: %v", err)
    }
}

func HandleListObjects(w http.ResponseWriter, r *http.Request) {
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
        log.Printf("Invalid max-keys value: %v", err)
        return
    }

    if bucketName == "" {
        http.Error(w, "Bucket name is required", http.StatusBadRequest)
        log.Println("Bucket name is missing")
        return
    }

    objects, err := storage.ListObjects(bucketName, prefix, marker, maxKeysInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error while listing objects: %v", err)
        return
    }

    w.Header().Set("Content-Type", "application/xml")
    w.WriteHeader(http.StatusOK)
    if err := xml.NewEncoder(w).Encode(objects); err != nil {
        http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
        log.Printf("Error while encoding XML: %v", err)
    }
}

func HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request: %s %s", r.Method, r.URL.Path)
    vars := mux.Vars(r)
    bucketName := vars["bucketName"]

    if bucketName == "" {
        http.Error(w, "Bucket name is required", http.StatusBadRequest)
        log.Println("Bucket name is missing")
        return
    }

    err := storage.DeleteBucket(bucketName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error while deleting bucket: %v", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
    log.Printf("Bucket %s successfully deleted", bucketName)
}

func HandleDeleteObject(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received DELETE request: %s %s", r.Method, r.URL.Path)
    vars := mux.Vars(r)
    bucketName := vars["bucketName"]
    objectName := vars["objectName"]

    if bucketName == "" || objectName == "" {
        http.Error(w, "Bucket name and object name are required", http.StatusBadRequest)
        log.Println("Bucket name or object name is missing")
        return
    }

    err := storage.DeleteObject(bucketName, objectName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Error deleting object %s in bucket %s: %v", objectName, bucketName, err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Object successfully deleted"))
    log.Printf("Object %s in bucket %s successfully deleted", objectName, bucketName)
}


func HandleDeleteBatch(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received POST ?delete request for batch deletion: %s %s", r.Method, r.URL.Path)

    vars := mux.Vars(r)
    bucketName := vars["bucketName"]
    log.Printf("Bucket name: %s", bucketName)

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusInternalServerError)
        log.Printf("Error reading request body: %v", err)
        return
    }
    log.Printf("Request body: %s", string(body))

    var deleteReq dto.DeleteBatchRequest
    err = xml.Unmarshal(body, &deleteReq)
    if err != nil {
        http.Error(w, "Error parsing XML", http.StatusInternalServerError)
        log.Printf("Error parsing XML: %v", err)
        return
    }

    for _, objectToDelete := range deleteReq.Objects {
        log.Printf("Attempting to delete object: %s", objectToDelete.Key)
        err := storage.DeleteObject(bucketName, objectToDelete.Key)
        if err != nil {
            http.Error(w, "Error deleting object", http.StatusInternalServerError)
            log.Printf("Error deleting object %s: %v", objectToDelete.Key, err)
            return
        }
        log.Printf("Successfully deleted object: %s", objectToDelete.Key)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Batch delete successful"))
}
