package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}

	log.Printf("Config loaded from %s", configPath)
	log.Printf("Database: %s:%d/%s", AppConfig.Database.Host, AppConfig.Database.Port, AppConfig.Database.Database)

	return nil
}