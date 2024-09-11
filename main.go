package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "my-s3-clone/handlers"
    "my-s3-clone/middleware"
    "os"
)

func init() {
    if _, err := os.Stat("./buckets"); os.IsNotExist(err) {
        log.Printf("Le répertoire 'buckets' n'existe pas. Création...")
        if err := os.Mkdir("./buckets", os.ModePerm); err != nil {
            log.Fatalf("Erreur lors de la création du répertoire 'buckets': %v", err)
        }
    }
}

func main() {
    r := mux.NewRouter()
    r.Use(middleware.LogRequestMiddleware)
    r.Use(middleware.LogResponseMiddleware) 

    // Route pour le health check
    r.HandleFunc("/probe-bsign{suffix:.*}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("<Response></Response>"))  
    }).Methods("GET", "HEAD")

    // Route spécifique pour la suppression par lot d'objets
    r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBatch).Queries("delete", "").Methods("POST")

    // Routes spécifiques aux objets
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleDeleteObject).Methods("DELETE")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleAddObject).Methods("POST", "PUT")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleCheckObjectExist).Methods("HEAD")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleDownloadObject).Methods("GET")

    // Routes spécifiques aux buckets
    r.HandleFunc("/{bucketName}/", handlers.HandleListObjects).Methods("GET", "HEAD")
    r.HandleFunc("/{bucketName}/", handlers.HandleGetBucket).Methods("GET")
    r.HandleFunc("/{bucketName}/", handlers.HandleCreateBucket).Methods("PUT", "GET", "HEAD")
    r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBucket).Methods("DELETE")

    // Route pour lister tous les buckets
    r.HandleFunc("/", handlers.HandleListBuckets).Methods("GET")

    log.Println("Serving on :9090")
    log.Fatal(http.ListenAndServe(":9090", r))
}
