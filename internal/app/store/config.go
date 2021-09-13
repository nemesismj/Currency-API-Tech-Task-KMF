package store

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// Config - store config struct
type Config struct {
	DatabaseUrl string `json:"databaseUrl"`
}

// pathConst
const pathConst = "./configs/"

// NewConfig - object of config
func NewConfig() *Config {
	config, _ := LoadConfiguration(pathConst + "config.json")
	return &Config{
		DatabaseUrl: config.DatabaseUrl,
	}
}
// LoadConfiguration - json object parser
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		logrus.Error(fmt.Sprintf("Произошла ошибка при загрузке конфигураций %s", err.Error()))
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
