package model

import (
	"encoding/xml"
	"jt_converter/internal/lib/random"
	"time"
)

type PMI struct {
	XMLName  xml.Name               `xml:"Label" json:"-"`
	Type     string                 `json:"type"`
	UID      string                 `json:"uid"`
	Name     string                 `xml:"value,attr"  json:"-"`
	RawProps Attributes             `xml:"Properties" json:"-"`
	Props    map[string]AWCProperty `json:"props"`
}

type AWCProperty struct {
	Type    string `json:"type"`
	UiValue string `json:"uiValue"`
	Value   string `json:"value"`
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

func (p *PMI) BuildAttributes() {
	p.UID = random.NewRandomString(15)
	p.Type = "xml_pmi"
	if p.Props == nil {
		p.Props = make(map[string]AWCProperty, 10)
	}
	p.Props["name"] = AWCProperty{
		Type:    "STRING",
		UiValue: p.Name,
		Value:   p.Name,
	}
	for _, propInstance := range p.RawProps.PropertyList {
		p.Props[propInstance.Key] = AWCProperty{
			Type:    "STRING",
			UiValue: propInstance.Value,
			Value:   propInstance.Value,
		}
	}
	if pmiType, ok := p.Props["NX_PMI_TYPE"]; ok && pmiType.Value == "22" {
		p.Type = "xml_pmi_radial"
	}
}
