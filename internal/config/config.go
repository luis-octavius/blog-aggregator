package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
	configFileMaxSize = 200 
	defaultUser = "default"
)

// Config represents application configuration settings 
type Config struct {
	Db_url            string `json:"db_url"`              // database connection URL 
	Current_user_name string `json:"current_user_name"`   // currently authenticated user 
}

// Read loads and parses the configuration from the config file.
// it searches for the config file in standard locations and falls back 
// to default values if the file doesn't exist or contains invalid data. 
// returns the configuration or an error if the file exists but can't be read. 
func Read() (Config, error) {
	// locate config file in filesystem
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// open config file for reading 
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %v\n", err)
	}

	defer configFile.Close()

	// read file contents into buffer 
	data := make([]byte, configFileMaxSize)
	_, err = configFile.Read(data)
	if err != nil {
		return Config{}, fmt.Errorf("error reading the contents of config file: %v\n", err)
	}

	// remove null bytes from buffer 
	trimmedBytes := bytes.Trim(data, "\x00")

	// default configuration used as fallback 
	// these values are used if config file doesn't exist or JSON parsing fails 
	cfg := Config{
		Db_url:            "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
		Current_user_name: defaultUser,
	}

	// parse JSON into config struct, overriding defaults 
	// If JSON is invalid, defaults are preserved and error is returned 
	err = json.Unmarshal(trimmedBytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling data from config file: %v\n", err)
	}

	return cfg, nil
}

// SetUser updates the current user in the configuration and persists it to disk.
func (cfg Config) SetUser(currentUser string) error {
	cfg.Current_user_name = currentUser

	// persists changes to config file 
	err := write(cfg)
	if err != nil {
		return fmt.Errorf("error setting user: %v\n", err)
	}

	return nil
}
