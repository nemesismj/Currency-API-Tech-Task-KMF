package apiserver

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"techtask/internal/app/store"
)

// pathConst
const pathConst = "./configs/"

// Config struct
type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	LogLevel string `json:"logLevel"`
	Store    *store.Config
}
// NewConfig - new api server config
func NewConfig() *Config {
	config, _ := LoadConfiguration(pathConst + "config.json")
	return &Config{
		Host:     config.Host,
		Port:     config.Port,
		LogLevel: config.LogLevel,
		Store:    store.NewConfig(),
	}
}
// LoadConfiguration - json object parser
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
