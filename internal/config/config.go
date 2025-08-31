package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() Config {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	jsonCfg, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("couldn't read the config file at ~/.gatorconfig.json")
	}

	var config Config

	json.Unmarshal(jsonCfg, &config)

	return config
}

func SetUser(username string) {
	if err := write(username); err != nil {
		log.Fatal("couldn't write jsonConfig with user set")
	}
}

func write(username string) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	config := Read()
	config.CurrentUsername = username

	jsonCfg, err := json.Marshal(&config)
	if err != nil {
		log.Fatal("couldn't set username.")
	}

	return os.WriteFile(configFilePath, jsonCfg, 0666)
}

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("couldn't get the user home directory")
	}

	jsonFilePath := fmt.Sprintf("%s/%s", userHomeDir, configFileName)
	_, err = os.Stat(jsonFilePath)
	if err != nil {
		return "", errors.New("couldn't find the config in the user home directory")
	}

	return jsonFilePath, nil
}
