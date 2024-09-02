package storage

import (
	"encoding/xml"
	"sync"
)

type Bucket struct {
	XMLName xml.Name           `xml:"Bucket"`
	Name    string             `xml:"Name"`
	Objects map[string]*Object `xml:"Objects>Object"`
	mu      sync.Mutex
}

// Création d'un nouveau bucket 
func NewBucket(name string) *Bucket {
	return &Bucket{
		Name:    name,
		Objects: make(map[string]*Object),
	}
}

// Conversion en XML
func (b *Bucket) ToXML() ([]byte, error) {
	return xml.MarshalIndent(b, "", "  ")
}
