package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %v\n", err)
	}

	defer configFile.Close()

	data := make([]byte, 200)
	_, err = configFile.Read(data)
	if err != nil {
		return Config{}, fmt.Errorf("error reading the contents of config file: %v\n", err)
	}

	trimmedBytes := bytes.Trim(data, "\x00")

	cfg := Config{
		Db_url:            "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
		Current_user_name: "default",
	}

	err = json.Unmarshal(trimmedBytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling data from config file: %v\n", err)
	}

	return cfg, nil
}

func (cfg Config) SetUser(currentUser string) error {
	cfg.Current_user_name = currentUser

	err := write(cfg)
	if err != nil {
		fmt.Printf("error setting user: %v\n", err)
	}

	return nil
}
