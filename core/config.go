package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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
	Elastic struct {
		Uris                      []string      `json:"uris"`
		Username                  string        `json:"username"`
		Password                  string        `json:"password"`
		SetSniff                  bool          `json:"db"`
		SnifferTimeout            time.Duration `json:"sniffer_time_out"`
		SnifferTimeoutStartup     time.Duration `json:"sniffer_time_out_startup"`
		SnifferInterval           time.Duration `json:"sniffer_interval"`
		SetHealthcheck            bool          `json:"set_healthcheck"`
		HealthcheckTimeout        time.Duration `json:"healthcheck_time_out"`
		HealthcheckTimeoutStartup time.Duration `json:"healthcheck_time_out_startup"`
		HealthcheckInterval       time.Duration `json:"healthcheck_interval"`
		RequiredPlugins           []string      `json:"plugins"`
		SetGzip                   bool          `json:"¨set_gzip"`
	} `json:"elastic"`
}

// NewConfig initiates new config instance
func NewConfig(configPath string) *Config {
	config := loadConfiguration(configPath)
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
