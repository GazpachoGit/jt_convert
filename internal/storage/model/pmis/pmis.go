package model

import (
	"encoding/xml"
	"time"
)

type PMI struct {
	XMLName    xml.Name `xml:"Label"`
	Type       string
	Name       string     `xml:"value,attr"`
	Attributes Attributes `xml:"Properties"`
}

type Attributes struct {
	PropertyList []Attribute `xml:"Property"`
}

type Attribute struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type Model struct {
	CreationDate time.Time
	JTFileName   string
	PMIs         []PMI
}

func NewModel(jtFileName string, pmis []PMI) *Model {
	return &Model{
		CreationDate: time.Now(),
		JTFileName:   jtFileName,
		PMIs:         pmis,
	}
}
