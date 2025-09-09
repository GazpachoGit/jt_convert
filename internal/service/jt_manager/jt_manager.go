package jtmng

import (
	"fmt"
	model "jt_converter/internal/storage/model/pmis"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const JTFormat = ".jt"
const XMLFormat = ".xml"

type Storage interface {
	SavePMIs(key string, m *model.Model) error
	GetPMIs(keys []string) ([]*model.Model, error)
	GetKeysList() ([]string, error)
}

type XMLManager interface {
	ParsePMIsFromXML(xmlFileFullPath string) ([]model.PMI, error)
}

type TCService interface {
	LoadFile(uid, typeName, name string) error
}

type JTManager struct {
	visualizerPath string
	xmlStoragePath string
	jtStoragePath  string
	log            *slog.Logger
	strg           Storage
	x              XMLManager
	tcService      TCService
}

func New(visualizerPath string, jtStoragePath string, xmlStoragePath string, strg Storage, x XMLManager, tc TCService, log *slog.Logger) *JTManager {
	log.Debug("create JT manager instance",
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
		x:              x,
		tcService:      tc,
	}
}

func (jt *JTManager) LoadFile(uid, typeName, name string) error {
	const op = "service.jtconverter.LoadJTFileFromTCWithPMIs"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", uid))
	//load file from TC using the ticket
	err := jt.tcService.LoadFile(uid, typeName, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("JT file saved")
	//process PMIs
	_, err = jt.GetPMIs(name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (jt *JTManager) GetPMIs(jtFileName string) (*model.Model, error) {
	const op = "service.jtconverter.GetPMIs"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))
	//check db
	models, err := jt.strg.GetPMIs([]string{jtFileName})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(models) > 0 && models[0] != nil {
		log.Info("Model found in DB!")
		return models[0], nil
	}
	log.Debug("Model not found in DB. Generating new from the JT...")
	//check jt file
	jtFilesList, err := jt.GetJTList()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !slices.Contains(jtFilesList, jtFileName) {
		return nil, fmt.Errorf("%s: %s", op, "Can't find the JT file in the JT directory")
	}
	//check xml
	xmlFilesList, err := jt.getXMLList()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !slices.Contains(xmlFilesList, jtFileName) {
		//convert
		err = jt.ConvertJTtoXML(jtFileName)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		log.Debug("XML file detected. Skip conversion from JT")
	}
	//parse pmis
	model, err := jt.ParsePMIsFromXML(jtFileName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//store pmis
	_, err = jt.StorePMIsInDB(model)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Successfully finished GetPMIs request!")
	return model, nil
}

func (jt *JTManager) GetJTList() ([]string, error) {
	const op = "service.jtconverter.GetJTList"
	log := jt.log.With(slog.String("op", op))

	log.Debug("searching JTs in folder", slog.String("dir", jt.jtStoragePath))
	files, err := os.ReadDir(jt.jtStoragePath)
	if err != nil {
		err = fmt.Errorf("%s: %w", "Error reading JT storage dir directory", err)
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

func (jt *JTManager) getXMLList() ([]string, error) {
	const op = "service.jtconverter.getXMLList"
	log := jt.log.With(slog.String("op", op))

	log.Debug("searching XMLs in folder", slog.String("dir", jt.xmlStoragePath))
	files, err := os.ReadDir(jt.xmlStoragePath)
	if err != nil {
		err = fmt.Errorf("%s: %w", "Error reading XML storage dir directory", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resp := make([]string, 0, 10)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), XMLFormat) {
			base := strings.TrimSuffix(file.Name(), XMLFormat)
			resp = append(resp, base)
		}
	}
	log.Info("found XML files: ", slog.Int("total", len(resp)))
	return resp, nil
}

func (jt *JTManager) ConvertJTtoXML(jtFileName string) error {
	const op = "service.jtconverter.ConvertJTtoXML"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))

	jtFileFullPath := jt.jtStoragePath + "\\" + jtFileName + JTFormat
	xmlOutputFileFullPath := jt.xmlStoragePath + "\\" + jtFileName + XMLFormat

	cmd := exec.Command(jt.visualizerPath, jtFileFullPath, "-pmi", "-output", xmlOutputFileFullPath)

	log.Debug("Executing command", slog.String("cmd", cmd.String()))

	//It overrides the file
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Cmd run without errors", slog.String("output", string(output)))

	if _, err := os.Stat(xmlOutputFileFullPath); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("XML file created", slog.String("xmlFilePath", string(xmlOutputFileFullPath)))
	return nil
}

func (jt *JTManager) ParsePMIsFromXML(jtFileName string) (*model.Model, error) {
	const op = "service.jtconverter.ParsePMIsFromXML"
	log := jt.log.With(slog.String("op", op), slog.String("jtFileName", jtFileName))
	log.Debug("start parse PMIs in the XML")
	xmlOutputFileFullPath := jt.xmlStoragePath + "\\" + jtFileName + ".xml"
	pmis, err := jt.x.ParsePMIsFromXML(xmlOutputFileFullPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	model := model.NewModel(jtFileName, pmis)
	log.Debug("successfully retrieved PMIs")
	return model, nil
}

func (jt *JTManager) StorePMIsInDB(m *model.Model) (string, error) {
	const op = "service.StorePMIsInDB"
	log := jt.log.With(slog.String("op", op))

	log.Debug("start StorePMIsInDB")
	if m.JTFileName == "" {
		return "", fmt.Errorf("%s: %s", "empty JTFileName detected", op)
	}
	err := jt.strg.SavePMIs(m.JTFileName, m)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("successfully finish StorePMIsInDB")
	return m.JTFileName, nil
}

func (jt *JTManager) GetPMIsList() ([]string, error) {
	const op = "service.jtconverter.GetPMIsList"
	log := jt.log.With(slog.String("op", op))

	log.Debug("start GetPMIsList")
	pmis, err := jt.strg.GetKeysList()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("found JT files: ", slog.Int("total", len(pmis)))
	return pmis, nil
}
