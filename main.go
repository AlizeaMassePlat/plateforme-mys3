package main

import (
	"fmt"
	"net/http"
	"strings"
	"my-s3-clone/storage"
)

func main() {
	// Initialisation du client MinIO
	storage.InitMinioClient("minio:9000", "minioadmin", "minioadmin")

	// Gestionnaire des routes HTTP
	http.HandleFunc("/", handleRequest)

	// Lance le serveur HTTP
	fmt.Println("Serveur en cours d'exécution sur :9090")
	http.ListenAndServe(":9090", nil)
}

// Gestion des requêtes HTTP
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Vérifie si c'est une requête de sondage
	if strings.Contains(r.URL.Path, "/probe-") {
		ProbeHandler(w, r)
		return
	}

	// Extraire le chemin de la requête
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 1 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	bucketName := pathParts[0]

	switch r.Method {
	case http.MethodPut:
		// Gestion de la création de bucket
		if len(pathParts) == 1 {
			handleCreateBucket(w, r, bucketName)
		} else {
			// Gestion des objets si besoin
			http.Error(w, "Object operations not supported yet", http.StatusNotImplemented)
		}
	case http.MethodGet:
		// Gestion de la liste des buckets
		handleListBuckets(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Fonction pour gérer les requêtes de sondage (probes)
func ProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Probe request received:", r.URL.Path)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<ProbeResult>OK</ProbeResult>")
}

// Gestion de la création d'un bucket
func handleCreateBucket(w http.ResponseWriter, r *http.Request, bucketName string) {
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	// Vérifier si le bucket existe déjà
	exists, err := storage.BucketExists(bucketName)
	if err != nil {
		http.Error(w, "Error checking bucket existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Bucket already exists", http.StatusConflict)
		return
	}

	err = storage.CreateBucket(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK) // Utilisez 200 OK après la création réussie
}

// Gestion de la liste des buckets
func handleListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := storage.ListBuckets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	xmlData, _ := buckets.ToXML()
	w.Write(xmlData)
}
