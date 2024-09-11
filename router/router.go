package router

import (
    "github.com/gorilla/mux"
    "my-s3-clone/handlers"
    "my-s3-clone/middleware"
    "my-s3-clone/storage"
    "net/http"
)

// SetupRouter sets up the router with default storage
func SetupRouter() *mux.Router {
    // Use default storage (e.g., FileStorage)
    return SetupRouterWithStorage(&storage.FileStorage{})
}

// SetupRouterWithStorage allows injecting custom storage (e.g., mock storage for tests)
func SetupRouterWithStorage(s storage.Storage) *mux.Router {
    r := mux.NewRouter()
    r.Use(middleware.LogRequestMiddleware)
    r.Use(middleware.LogResponseMiddleware)

    // Health check route
    r.HandleFunc("/probe-bsign{suffix:.*}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("<Response></Response>"))
    }).Methods("GET", "HEAD")

    // Batch delete route
    r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBatch(s)).Queries("delete", "").Methods("POST")

    // Object-specific routes
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleDeleteObject(s)).Methods("DELETE")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleAddObject(s)).Methods("POST", "PUT")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleCheckObjectExist(s)).Methods("HEAD")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleDownloadObject(s)).Methods("GET")

    // Bucket-specific routes
    r.HandleFunc("/{bucketName}/", handlers.HandleListObjects(s)).Methods("GET", "HEAD")
    r.HandleFunc("/{bucketName}/", handlers.HandleGetBucket(s)).Methods("GET")
    r.HandleFunc("/{bucketName}/", handlers.HandleCreateBucket(s)).Methods("PUT", "GET", "HEAD")
    r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBucket(s)).Methods("DELETE")

    // Route for listing all buckets
    r.HandleFunc("/", handlers.HandleListBuckets(s)).Methods("GET")

    return r
}
