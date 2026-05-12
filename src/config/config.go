package config

import (
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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file %s not found, using environment variables", configPath)
		} else {
			return err
		}
	} else {
		err = yaml.Unmarshal(data, &AppConfig)
		if err != nil {
			return err
		}
		log.Printf("Config loaded from %s", configPath)
	}

	overrideFromEnv()

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
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		AppConfig.JWT.Secret = secret
	}
	if expireHour := os.Getenv("JWT_EXPIRE_HOUR"); expireHour != "" {
		if h, err := strconv.Atoi(expireHour); err == nil {
			AppConfig.JWT.ExpireHour = h
		}
	}
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		if p, err := strconv.Atoi(serverPort); err == nil {
			AppConfig.Server.Port = p
		}
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		AppConfig.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		if p, err := strconv.Atoi(redisPort); err == nil {
			AppConfig.Redis.Port = p
		}
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		AppConfig.Redis.Password = redisPassword
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		if db, err := strconv.Atoi(redisDB); err == nil {
			AppConfig.Redis.DB = db
		}
	}
}
