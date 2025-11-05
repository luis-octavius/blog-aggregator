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

func Read() Config {
	configFilePath, err := getConfigFilePath()
	fmt.Println("config file path: ", configFilePath)
	if err != nil {
		return Config{}
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("error opening config file: %v\n", err)
		return Config{}
	}

	defer configFile.Close()

	data := make([]byte, 200)
	_, err = configFile.Read(data)
	if err != nil {
		fmt.Printf("error reading the contents of config file: %v\n", err)
	}

	trimmedBytes := bytes.Trim(data, "\x00") //

	cfg := Config{
		Db_url:            "postgres://example",
		Current_user_name: "default",
	}

	err = json.Unmarshal(trimmedBytes, &cfg)
	if err != nil {
		fmt.Printf("error unmarshaling data from config file: %v\n", err)
		return Config{}
	}

	return cfg
}

func (cfg Config) SetUser(currentUser string) error {
	cfg.Current_user_name = currentUser

	err := write(cfg)
	if err != nil {
		fmt.Printf("error setting user: %v\n", err)
	}

	return nil
}
