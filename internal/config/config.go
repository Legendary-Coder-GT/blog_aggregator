package config

import (
	"encoding/json"
	"os"
)

const configFileName string = "/.gatorconfig.json"

type Config struct {
	DB_URL string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(jsonData, c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

func (c *Config) SetUser(new_user string) error {
	c.Current_user_name = new_user
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(*cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := homeDir + configFileName
	return path, nil
}