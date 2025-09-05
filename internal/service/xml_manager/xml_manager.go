package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	model "jt_converter/internal/storage/model/pmis"
	"log/slog"
	"os"
)

type XMLManager struct {
	log *slog.Logger
}

func NewXMLManager(log *slog.Logger) *XMLManager {
	return &XMLManager{log}
}

func (x *XMLManager) ParsePMIsFromXML(xmlFileFullPath string) ([]model.PMI, error) {
	const op = "service.xml_Manager.ParsePMIsFromXML"
	log := x.log.With(slog.String("op", op), slog.String("file", xmlFileFullPath))

	file, err := os.Open(xmlFileFullPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	pmis := make([]model.PMI, 0, 10)

	for {
		token, err := decoder.Token()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Debug("Error during xml file reading", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Label" {
				var label model.PMI
				err := decoder.DecodeElement(&label, &se)
				if err != nil {
					log.Debug("err during interactive decoding", slog.String("err", err.Error()))
					continue
				}
				log.Debug("Label detected", slog.String("labelName", label.Name))
				pmis = append(pmis, label)
			}
		}
	}
	log.Info("total PMIs found", slog.Int("total", len(pmis)))
	return pmis, nil
}
