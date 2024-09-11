package storage

import (
	"io"
	"time"
	"my-s3-clone/dto"
	"os"

)

// Storage interface définissant les méthodes de gestion des objets et des buckets
type Storage interface {
    AddObject(bucketName, objectName string, data io.Reader, contentSha256 string) error
    DeleteObject(bucketName, objectName string) error
    DeleteBucket(bucketName string) error
    GetObject(bucketName, objectName string) ([]byte, FileInfo, error)
    CheckObjectExist(bucketName, objectName string) (bool, time.Time, int64, error)
    CheckBucketExists(bucketName string) (bool, error)
    ListBuckets() []string
    ListObjects(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error)
    CreateBucket(bucketName string) error
}

// FileInfo représente les métadonnées d'un fichier (objet)
type FileInfo interface {
    Name() string       // Nom de base du fichier
    Size() int64        // Taille logique du fichier en octets
    Mode() os.FileMode  // Informations sur le mode de fichier
    ModTime() time.Time // Heure de dernière modification
    IsDir() bool        // Indique si c'est un répertoire
    Sys() interface{}   // Données spécifiques au système sous-jacent
}
