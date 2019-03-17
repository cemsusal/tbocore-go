package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config configures the application
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     string `json:"port"`
		User     string `json:"user"`
		DbName   string `json:"dbname"`
		SslMode  string `json:"sslmode"`
	} `json:"database"`
	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
}

// NewConfig initiates new config instance
func NewConfig(configFilePath string) *Config {
	config := loadConfiguration(configFilePath)
	return &config
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
