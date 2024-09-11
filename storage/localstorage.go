package storage

import (
    "strings"
    "os"
    "path/filepath"
    "log"
    "fmt"
    "io"
    "bufio"  
    "strconv"
    "time"
    "my-s3-clone/dto"
)

// FileStorage implémente l'interface Storage avec un stockage basé sur le système de fichiers
type FileStorage struct{}

const storageRoot = "/mydata/data"

func ProcessChunkedStream(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)
    log.Println("Started processing chunked stream")
    
    for {
        line, err := bufReader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading chunk size: %v", err)
            return fmt.Errorf("error reading chunk size: %v", err)
        }
        log.Printf("Received chunk size line: %s", line)

        line = strings.TrimSpace(line)
        parts := strings.SplitN(line, ";", 2)
        chunkSizeHex := parts[0]

        chunkSize, err := strconv.ParseInt(chunkSizeHex, 16, 64)
        if err != nil {
            log.Printf("Error parsing chunk size: %v", err)
            return fmt.Errorf("error parsing chunk size: %v", err)
        }

        log.Printf("Parsed chunk size: %d", chunkSize)

        if chunkSize == 0 {
            log.Println("Received final chunk (size 0), finishing")
            break
        }

        if _, err := io.CopyN(writer, bufReader, chunkSize); err != nil {
            log.Printf("Error reading chunk data: %v", err)
            return fmt.Errorf("error reading chunk data: %v", err)
        }

        if _, err := bufReader.Discard(2); err != nil {
            log.Printf("Error discarding CRLF: %v", err)
            return fmt.Errorf("error discarding CRLF: %v", err)
        }

        log.Println("Successfully processed a chunk")

        if len(parts) > 1 {
            chunkSignature := parts[1]
            log.Printf("Chunk signature: %s", chunkSignature)
        }
    }

    log.Println("Completed processing chunked stream")
    return nil
}

// Ajout d'un objet dans un bucket
func (fs *FileStorage) AddObject(bucketName, objectName string, data io.Reader, contentSha256 string) error {
    objectPath, err := getUniqueObjectPath(bucketName, objectName)
    if err != nil {
        log.Printf("Failed to create object path: %v", err)
        return fmt.Errorf("Failed to create object path: %v", err)
    }

    file, err := os.Create(objectPath)
    if err != nil {
        log.Printf("Failed to create file: %v", err)
        return fmt.Errorf("Failed to create file: %v", err)
    }
    defer file.Close()

    if err := writeObjectToFile(data, file, contentSha256); err != nil {
        return err
    }

    log.Printf("Successfully uploaded file: %s", objectPath)
    return nil
}

// Fonction pour obtenir un chemin unique si l'objet existe déjà
func getUniqueObjectPath(bucketName, objectName string) (string, error) {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)
    if _, err := os.Stat(objectPath); os.IsNotExist(err) {
        return objectPath, nil
    }

    log.Printf("Object already exists, generating new name for: %s", objectPath)

    objectNameWithoutExt := strings.TrimSuffix(objectName, filepath.Ext(objectName))
    extension := filepath.Ext(objectName)
    newObjectName := objectNameWithoutExt
    suffix := 1

    for {
        newObjectName = fmt.Sprintf("%s-%d%s", objectNameWithoutExt, suffix, extension)
        newObjectPath := filepath.Join(storageRoot, bucketName, newObjectName)
        if _, err := os.Stat(newObjectPath); os.IsNotExist(err) {
            return newObjectPath, nil
        }
        suffix++
    }
}

// Fonction qui gère l'écriture du flux dans le fichier
func writeObjectToFile(data io.Reader, file *os.File, contentSha256 string) error {
    if contentSha256 == "STREAMING-AWS4-HMAC-SHA256-PAYLOAD" {
        log.Println("Processing as chunked stream")
        if err := ProcessChunkedStream(data, file); err != nil {
            log.Printf("Failed to write chunked data: %v", err)
            return fmt.Errorf("Failed to write chunked data: %v", err)
        }
    } else {
        log.Println("Processing as regular stream")
        if _, err := io.Copy(file, data); err != nil {
            log.Printf("Failed to write data: %v", err)
            return fmt.Errorf("Failed to write data: %v", err)
        }
    }
    return nil
}

// Lister les objets dans un bucket
func (fs *FileStorage) ListObjects(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error) {
    bucketPath := filepath.Join(storageRoot, bucketName)

    objects, err := filepath.Glob(filepath.Join(bucketPath, prefix+"*"))
    if err != nil {
        return dto.ListObjectsResponse{}, fmt.Errorf("error while listing objects: %v", err)
    }

    response := dto.ListObjectsResponse{
        Xmlns:       "http://s3.amazonaws.com/doc/2006-03-01/",
        Name:        bucketName,
        Prefix:      prefix,
        Marker:      marker,
        MaxKeys:     maxKeys,
        IsTruncated: false,
        Contents:    make([]dto.Object, 0),
    }

    for i, object := range objects {
        if i >= maxKeys {
            response.IsTruncated = true
            break
        }

        fileInfo, err := os.Stat(object)
        if err != nil {
            return dto.ListObjectsResponse{}, fmt.Errorf("error retrieving file info: %v", err)
        }

        response.Contents = append(response.Contents, dto.Object{
            Key:          filepath.Base(object),
            LastModified: fileInfo.ModTime(),
            Size:         int(fileInfo.Size()),
        })
    }

    return response, nil
}

// Lister les buckets
func (fs *FileStorage) ListBuckets() []string {
    var buckets []string
    files, err := os.ReadDir(storageRoot)
    if err != nil {
        return buckets
    }

    for _, file := range files {
        if file.IsDir() {
            buckets = append(buckets, file.Name())
        }
    }

    return buckets
}

// Créer un bucket
func (fs *FileStorage) CreateBucket(bucketName string) error {
    bucketPath := filepath.Join(storageRoot, bucketName)
    if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
        return err
    }
    return nil
}

// Récupération d'un objet dans un bucket
func (fs *FileStorage) GetObject(bucketName, objectName string) ([]byte, FileInfo, error) {
	objectPath := filepath.Join(storageRoot, bucketName, objectName)
	log.Printf("Tentative de récupération de l'objet : %s", objectPath)

	// Lire le fichier
	data, err := os.ReadFile(objectPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture de l'objet: %v", err)
		return nil, nil, err
	}

	// Récupérer les métadonnées du fichier
	fileInfo, err := os.Stat(objectPath)
	if err != nil {
		log.Printf("Erreur lors de la récupération des métadonnées du fichier: %v", err)
		return nil, nil, err
	}

	// Retourner le contenu du fichier et les métadonnées encapsulées dans fileInfoWrapper
	return data, &fileInfoWrapper{fileInfo: fileInfo}, nil
}

// Vérification de l'existence d'un objet dans un bucket
func (fs *FileStorage) CheckObjectExist(bucketName, objectName string) (bool, time.Time, int64, error) {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)

    fileInfo, err := os.Stat(objectPath)
    if os.IsNotExist(err) {
        return false, time.Time{}, 0, nil
    } else if err != nil {
        log.Printf("Error checking object: %v", err)
        return false, time.Time{}, 0, fmt.Errorf("error checking object existence: %v", err)
    }

    return true, fileInfo.ModTime(), fileInfo.Size(), nil
}

// Vérification de l'existence d'un bucket
func (fs *FileStorage) CheckBucketExists(bucketName string) (bool, error) {
    bucketPath := filepath.Join(storageRoot, bucketName)
    if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
        return false, nil
    } else if err != nil {
        return false, err
    }
    return true, nil
}

// Suppression d'un bucket
func (fs *FileStorage) DeleteBucket(bucketName string) error {
    bucketPath := filepath.Join(storageRoot, bucketName)

    if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
        log.Printf("Bucket %s does not exist", bucketName)
        return err
    }

    err := os.RemoveAll(bucketPath)
    if err != nil {
        log.Printf("Failed to delete bucket %s: %v", bucketName, err)
        return err
    }

    log.Printf("Bucket %s successfully deleted", bucketName)
    return nil
}

// Suppression d'un objet dans un bucket
func (fs *FileStorage) DeleteObject(bucketName, objectName string) error {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)

    if _, err := os.Stat(objectPath); os.IsNotExist(err) {
        log.Printf("Object %s does not exist in bucket %s", objectName, bucketName)
        return err
    }

    err := os.Remove(objectPath)
    if err != nil {
        log.Printf("Failed to delete object %s in bucket %s: %v", objectName, bucketName, err)
        return err
    }

    log.Printf("Object %s in bucket %s successfully deleted", objectName, bucketName)
    return nil
}