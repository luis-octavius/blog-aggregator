package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting the home path: %w", err)
	}

	configFilePath := path.Join(home, configFileName)
	return configFilePath, nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting the config file path: %w", err)
	}

	err = os.Remove(configFilePath)
	if err != nil {
		return fmt.Errorf("error removing the config file path: %w", err)
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("error creating the config file path: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to config file with json data: %w", err)
	}

	return nil
}
