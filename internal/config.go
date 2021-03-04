package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirSuffix             = "deezerdl"
	configDirPerms  os.FileMode = 0755
	configFile                  = "config.json"
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

// GetConfigDir returns the platform-specific configuration directory
func GetConfigDir() (string, error) {
	// get the platform-specific config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// append the application-specific suffix
	configDir := filepath.Join(userConfigDir, configDirSuffix)

	return configDir, nil
}

// CreateConfig creates the config dir and config file if it doesn't
// exist
func CreateConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	fmt.Println(configDir)

	// check if config dir exists
	if exists, err := FileExists(os.ExpandEnv(configDir)); err != nil {
		return err
	} else if !exists {
		os.MkdirAll(os.ExpandEnv(configDir), configDirPerms)
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

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	// open config file
	fullPath := filepath.Join(os.ExpandEnv(configDir), configFile)
	inFile, err := os.Open(fullPath)
	if err != nil {
		return nil, err
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
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
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
