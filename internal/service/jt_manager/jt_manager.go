package jtmng

import (
	"fmt"
	model "jt_converter/internal/storage/model/pmis"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

const JTFormat = ".jt"
const XMLFormat = ".xml"

type Storage interface {
	SavePMIs(key string, m *model.Model) error
	GetPMIs(keys []string) ([]*model.Model, error)
	GetKeysList() ([]string, error)
}

type JTManager struct {
	visualizerPath string
	xmlStoragePath string
	jtStoragePath  string
	log            *slog.Logger
	strg           Storage
}

func New(visualizerPath string, jtStoragePath string, xmlStoragePath string, strg Storage, log *slog.Logger) *JTManager {
	log.Info("create JT manager instance",
		slog.String("visualizerPath", visualizerPath),
		slog.String("jtStoragePath", jtStoragePath),
		slog.String("xmlStoragePath", xmlStoragePath),
	)
	return &JTManager{
		visualizerPath: visualizerPath + "\\JTInspector.exe",
		xmlStoragePath: xmlStoragePath,
		jtStoragePath:  jtStoragePath,
		log:            log,
		strg:           strg,
	}
}

func (jt *JTManager) GetPMIs(jtFileName string) (*model.Model, error) {
	const op = "service.jtconverter.GetPMIs"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))
	//check db
	models, err := jt.strg.GetPMIs([]string{jtFileName})
	if err != nil {
		log.Error("Error during GetPMIs storage request", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(models) > 0 && models[0] != nil {
		log.Info("Model found in DB!")
		return models[0], nil
	}
	log.Info("Model not found in DB. Generating new from the JT...")
	//check jt file
	jtFilesList, err := jt.GetJTList()
	if err != nil {
		log.Error("Can't get the list of JT files", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !slices.Contains(jtFilesList, jtFileName) {
		log.Info("Can't find the JT file in the JT directory")
		return nil, fmt.Errorf("%s: %s", op, "Can't find the JT file in the JT directory")
	}
	//convert
	err = jt.ConvertJTtoXML(jtFileName)
	if err != nil {
		log.Error("Can't convert JT to XML", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//parse pmis
	model, err := jt.ParsePMIsFromXML(jtFileName)
	if err != nil {
		log.Error("Can't parse PMIs from XML", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//store pmis
	_, err = jt.StorePMIsInDB(model)
	if err != nil {
		log.Error("Can't write pmis into DB", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Successfully procesed GetPMIs request!")
	return model, nil
}

func (jt *JTManager) GetJTList() ([]string, error) {
	const op = "service.jtconverter.GetJTList"
	log := jt.log.With(slog.String("op", op))

	log.Info("searching JTs in folder", slog.String("dir", jt.jtStoragePath))
	files, err := os.ReadDir(jt.jtStoragePath)
	if err != nil {
		log.Error("Error reading JT storage dir directory:", slog.String("error", err.Error()))
		fmt.Println("Error reading directory:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resp := make([]string, 0, 10)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), JTFormat) {
			base := strings.TrimSuffix(file.Name(), JTFormat)
			resp = append(resp, base)
		}
	}
	log.Info("found JT files: ", slog.Int("total", len(resp)))
	return resp, nil
}

func (jt *JTManager) ConvertJTtoXML(jtFileName string) error {
	const op = "service.jtconverter.ConvertJTtoXML"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))

	jtFileFullPath := jt.jtStoragePath + "\\" + jtFileName + JTFormat
	xmlOutputFileFullPath := jt.xmlStoragePath + "\\" + jtFileName + XMLFormat

	cmd := exec.Command(jt.visualizerPath, jtFileFullPath, "-pmi", jtFileFullPath, "-output", xmlOutputFileFullPath)

	log.Info("Executing command", slog.String("cmd", cmd.String()))

	//It overrides the file
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("failed to run the cmd", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Cmd run without errors", slog.String("output", string(output)))

	if _, err := os.Stat(xmlOutputFileFullPath); err == nil {
		log.Info("XML file created", slog.String("xmlFilePath", string(xmlOutputFileFullPath)))
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}

func (jt *JTManager) ParsePMIsFromXML(jtFileName string) (*model.Model, error) {
	const op = "service.jtconverter.ParsePMIsFromXML"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))
	log.Debug("start parse PMIs in the XML")
	//xmlOutputFileFullPath := jt.xmlStoragePath + "\\" + jtFileName + ".xml"
	return &model.Model{
		CreationDate: time.Now(),
		JTFileName:   jtFileName,
		PMIs: []model.PMI{
			{
				Label: "Diameter_1",
				Type:  "Dimension",
				Attributes: map[string]string{
					"value": "20",
					"upper": "0.05",
					"lower": "0.02",
				},
			},
		},
	}, nil
}

func (jt *JTManager) StorePMIsInDB(m *model.Model) (string, error) {
	const op = "service.StorePMIsInDB"
	log := jt.log.With(slog.String("op", op))

	log.Debug("start StorePMIsInDB")
	if m.JTFileName == "" {
		return "", fmt.Errorf("%s: %s", "can't get key from Model to store in DB ", op)
	}
	err := jt.strg.SavePMIs(m.JTFileName, m)
	if err != nil {
		log.Error("can't save the Model in db", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("successfully finish StorePMIsInDB")
	return m.JTFileName, nil
}

func (jt *JTManager) GetPMIsList() ([]string, error) {
	const op = "service.jtconverter.GetPMIsList"
	log := jt.log.With(slog.String("op", op))
	pmis, err := jt.strg.GetKeysList()
	if err != nil {
		log.Error("can't get list of PMIs", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("found JT files: ", slog.Int("total", len(pmis)))
	return pmis, nil
}
