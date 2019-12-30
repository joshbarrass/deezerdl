package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDir                  = "$HOME/.config/deezerdl"
	configDirPerms os.FileMode = 0755
	configFile                 = "config.json"
)

type Configuration struct {
	Version       string `json:"version"`
	ARLCookie     string `json:"arl"`
	DefaultFormat string `json:"default_format"`
}

// NewConfiguration creates an empty, default config
func NewConfiguration() *Configuration {
	return &Configuration{
		Version:       "1",
		DefaultFormat: "MP3_320",
	}
}

// CreateConfig creates the config dir and config file if it doesn't
// exist
func CreateConfig() error {
	// check if config dir exists
	if exists, err := FileExists(os.ExpandEnv(configDir)); err != nil {
		return err
	} else if !exists {
		os.Mkdir(os.ExpandEnv(configDir), configDirPerms)
	}

	// check if config file exists
	fullPath := filepath.Join(os.ExpandEnv(configDir), configFile)
	if exists, err := FileExists(fullPath); err != nil {
		return err
	} else if exists {
		// file exists, exit
		return nil
	}

	// create new default config and save
	config := NewConfiguration()
	config.SaveConfig()

	return nil
}

// LoadConfig loads the configuration file
func LoadConfig() (*Configuration, error) {
	// try to create a config if it doesn't exist
	if err := CreateConfig(); err != nil {
		return nil, err
	}

	// open config file
	fullPath := filepath.Join(os.ExpandEnv(configDir), configFile)
	inFile, err := os.Open(fullPath)
	if err != nil {
		return nil, nil
	}
	defer inFile.Close()

	// read data
	var config Configuration
	decoder := json.NewDecoder(inFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveConfig saves the config object. The configuration file must
// already exist -- call CreateConfig first if this is not the case.
func (config *Configuration) SaveConfig() error {
	// open config
	fullPath := filepath.Join(os.ExpandEnv(configDir), configFile)
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// write config to file
	encoder := json.NewEncoder(outFile)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
