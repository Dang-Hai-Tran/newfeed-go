package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LogLevel string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret           string
	ExpirationHours  int
	RefreshSecret    string
	RefreshDuration  time.Duration
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
