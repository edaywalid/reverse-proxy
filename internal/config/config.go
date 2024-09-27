package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	URL string `json:"url" yaml:"url"`
}

type Config struct {
	Servers   []ServerConfig `json:"servers" yaml:"servers"`
	CertFile  string         `json:"cert_file" yaml:"cert_file"`
	KeyFile   string         `json:"key_file" yaml:"key_file"`
	HTTPPort  int            `json:"http_port" yaml:"http_port"`
	HTTPSPort int            `json:"https_port" yaml:"https_port"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open file: %s", err))
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read file: %s", err))
	}

	var config Config
	ext := filepath.Ext(filePath)
	switch {
	case ext == ".json":
		err = json.Unmarshal(content, &config)
	case ext == ".yml" || ext == ".yaml":
		err = yaml.Unmarshal(content, &config)
	default:
		return nil, fmt.Errorf("unsupported file type: %v", filePath)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal config: %s", err))
	}

	return &config, nil
}
