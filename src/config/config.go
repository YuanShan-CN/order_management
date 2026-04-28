package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

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

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireHour int    `yaml:"expire_hour"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
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

	overrideFromEnv()

	log.Printf("Config loaded from %s", configPath)
	log.Printf("Database: %s:%d/%s", AppConfig.Database.Host, AppConfig.Database.Port, AppConfig.Database.Database)

	return nil
}

func overrideFromEnv() {
	if host := os.Getenv("DB_HOST"); host != "" {
		AppConfig.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			AppConfig.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		AppConfig.Database.Username = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		AppConfig.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		AppConfig.Database.Database = name
	}
}