package config

import (
	"log"
	"os" // 已弃用：Go 1.16+ 应使用 os 包中的 ReadFile
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

type SchedulerConfig struct {
	Timezone     string `yaml:"timezone"`
	ExportHour   int    `yaml:"export_hour"`
	ExportMinute int    `yaml:"export_minute"`
}

type Config struct {
	Database  DatabaseConfig  `yaml:"database"`
	Server    ServerConfig    `yaml:"server"`
	JWT       JWTConfig       `yaml:"jwt"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

var AppConfig Config

func LoadConfig(configPath string) error {
	// 设置默认值
	AppConfig.Scheduler.Timezone = "Asia/Shanghai"
	AppConfig.Scheduler.ExportHour = 0
	AppConfig.Scheduler.ExportMinute = 5

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
	log.Printf("Scheduler: timezone=%s, export time=%02d:%02d",
		AppConfig.Scheduler.Timezone,
		AppConfig.Scheduler.ExportHour,
		AppConfig.Scheduler.ExportMinute)

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
	if timezone := os.Getenv("TZ"); timezone != "" {
		AppConfig.Scheduler.Timezone = timezone
	}
	if exportHour := os.Getenv("EXPORT_HOUR"); exportHour != "" {
		if h, err := strconv.Atoi(exportHour); err == nil {
			AppConfig.Scheduler.ExportHour = h
		}
	}
	if exportMinute := os.Getenv("EXPORT_MINUTE"); exportMinute != "" {
		if m, err := strconv.Atoi(exportMinute); err == nil {
			AppConfig.Scheduler.ExportMinute = m
		}
	}
}
