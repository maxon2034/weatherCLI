package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DefaultCity string `yaml:"default_city"`
}

func Load() (Config, error) {
	conf := Config{}
	confFile := "C:/Dev/praxiscode/weatherCLI/internal/config/conf.json"
	_, err := os.Stat(confFile)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("Error in checking configuration file state: ", err)
	}
	confText, err := os.ReadFile(confFile)
	if err != nil {
		return Config{}, fmt.Errorf("Error in reading configuration file: ", err)
	}

	err = json.Unmarshal(confText, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("Error in deserializing json: ", err)
	}

	return conf, nil
}

func Save(cfg Config) error {
	confFile := "C:/Dev/praxiscode/weatherCLI/internal/config/conf.json"
	confText, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error in serealizing configuration: ", err)
	}

	err = os.WriteFile(confFile, confText, 0700)
	if err != nil {
		return fmt.Errorf("error in writing into configuration file: ", err)
	}

	return nil
}
