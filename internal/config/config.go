package config

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	Connections struct {
		Postgres struct {
			Host     string
			DB       string
			User     string
			Password string
		}
		Redis struct {
			Host         string
			Password     string
			WriteTimeout time.Duration
			ReadTimeout  time.Duration
		}
	}
	Swagger struct {
		Host     string
		BasePath string
	}
	Log struct {
		Level string
	}
}

func Parse() (Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config file failed: %w", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, fmt.Errorf("unmarshal to config struct failed: %w", err)
	}

	return c, nil
}
