package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// getConfigFilePath returns the full path of the configuration file 
// located in the user's home directory. 
// returns the absolute file path or an error if the home directory cannot be found. 
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting the home path: %w", err)
	}

	configFilePath := path.Join(home, configFileName)
	return configFilePath, nil
}

// write persists the configuration to disk as JSON. 
// it replaces the entire config file by removing the existing file 
// and creating a new one to ensure data consistency. 
// returns error if JSON marshaling fails or any file operation fails.
func write(cfg Config) error {
	// Marshal configuration to JSON format 
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	// determine config file location 
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting the config file path: %w", err)
	}

	// remove existing config file to ensure clean write 
	// ignore "file not found" errors as the file may not exist yet
	if err = os.Remove(configFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing the config file: %w", err)
	}

	// create new config file with appropriate permissions 
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("error creating the config file path: %w", err)
	}
	defer file.Close()

	// write JSON data to file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to config file with json data: %w", err)
	}

	return nil
}
