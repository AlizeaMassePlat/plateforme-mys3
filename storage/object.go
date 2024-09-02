package storage

import (
	"encoding/xml"
)

type Object struct {
	XMLName xml.Name `xml:"Object"`
	Name    string   `xml:"Name"`
	Data    []byte   `xml:"Data"`
}
