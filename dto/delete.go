package dto

import (
    "encoding/xml"
)
type DeleteObjectRequest struct {
	Quiet  bool `xml:"Quiet"`
	Object struct {
		Key string `xml:"Key"`
	} `xml:"Object"`
}

type DeleteResult struct {
	DeletedResult []Deleted `xml:"Deleted"`
}

type Deleted struct {
	Key string `xml:"Key"`
}

// DeleteBatchRequest représente la requête de suppression d'objets en batch
type DeleteBatchRequest struct {
    XMLName xml.Name         `xml:"Delete"`
    Objects []ObjectToDelete  `xml:"Object"`
}

// ObjectToDelete représente un objet à supprimer
type ObjectToDelete struct {
    Key string `xml:"Key"`
}
