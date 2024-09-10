package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "my-s3-clone/handlers"
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

func logRequestMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Requête reçue 1: %s %s", r.Method, r.RequestURI)

        if len(r.URL.Query()) > 0 {
            log.Printf("Query Params: %v", r.URL.Query())
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    r := mux.NewRouter()
    r.Use(logRequestMiddleware)

    r.HandleFunc("/probe-bsign{suffix:.*}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("<Response></Response>"))  
    }).Methods("GET", "HEAD")

    r.HandleFunc("/", handlers.HandleListBuckets).Methods("GET")

    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleAddObject).Methods("POST", "PUT")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleCheckObjectExist).Methods("HEAD")
    r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleDownloadObject).Methods("GET")

    r.HandleFunc("/{bucketName}/", handlers.HandleListObjects).Methods("GET", "HEAD")
    r.HandleFunc("/{bucketName}/", handlers.HandleGetBucket).Methods("GET")
    r.HandleFunc("/{bucketName}/", handlers.HandleCreateBucket).Methods("PUT", "GET", "HEAD")

    r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBucket).Methods("DELETE")

    log.Println("Serving on :9090")
    log.Fatal(http.ListenAndServe(":9090", r))
}
