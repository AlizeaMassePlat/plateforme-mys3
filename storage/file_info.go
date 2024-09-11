package storage

import (
	"os"
	"time"
)

// fileInfoWrapper encapsule un os.FileInfo pour implémenter l'interface FileInfo
type fileInfoWrapper struct {
	fileInfo os.FileInfo
}

// Implémentation des méthodes de l'interface FileInfo

func (fi *fileInfoWrapper) Name() string {
	return fi.fileInfo.Name()
}

func (fi *fileInfoWrapper) Size() int64 {
	return fi.fileInfo.Size()
}

func (fi *fileInfoWrapper) Mode() os.FileMode {
	return fi.fileInfo.Mode()
}

func (fi *fileInfoWrapper) ModTime() time.Time {
	return fi.fileInfo.ModTime()
}

func (fi *fileInfoWrapper) IsDir() bool {
	return fi.fileInfo.IsDir()
}

func (fi *fileInfoWrapper) Sys() interface{} {
	return fi.fileInfo.Sys()
}
