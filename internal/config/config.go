package config

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPSever ServerConfig `yaml:"http_server"`
	JT        JTConfig     `yaml:"jt"`
	TC        TCConfig     `yaml:"tc"`
}

type JTConfig struct {
	VisualizerPath string `yaml:"visualizer_path" env-requered:"true"`
	XmlStoragePath string `yaml:"xml_storage_path" env-requered:"true"`
	JtStoragePath  string `yaml:"js_storage_path" env-requered:"true"`
	DBPath         string `yaml:"db_path" env-requered:"true"`
}

type ServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:9000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type TCConfig struct {
	TCURL    string `yaml:"tc_url" env-default:"http://localhost:3000"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(err)
	}

	if _, err := os.Stat(cfg.JT.VisualizerPath); err != nil {
		panic("location does not exist at path: " + cfg.JT.VisualizerPath)
	}
	if _, err := os.Stat(cfg.JT.JtStoragePath); err != nil {
		panic("location does not exist at path: " + cfg.JT.JtStoragePath)
	}
	if _, err := os.Stat(cfg.JT.XmlStoragePath); err != nil {
		panic("location does not exist at path: " + cfg.JT.XmlStoragePath)
	}
	return cfg
}

// check flag or env var
func fetchConfigPath() string {
	var targetPath string
	flag.StringVar(&targetPath, "config", "", "Path to config file")
	flag.Parse()

	if targetPath == "" {
		targetPath = os.Getenv("CONFIG_PATH")
	}

	if targetPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exeDir := filepath.Dir(exePath)
		targetPath = filepath.Join(exeDir, "local.yml")
		if _, err := os.Stat(targetPath); err != nil {
			panic(err)
		}
	}
	return targetPath
}
