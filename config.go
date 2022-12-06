package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

const (
	configPath = "./configs/"
	configType = "yaml"
	configName = "config"
)

func init() {
	flag.String("environment", "", "Values: local, dev, prod")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Drivername string `yaml:"drivername"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Dbname     string `yaml:"dbname"`
	Sslmode    string `yaml:"sslmode"`
}

type jwtConfig struct {
	AccessSecret  string `yaml:"accesssecret"`
	RefreshSecret string `yaml:"refreshsecret"`
}

type applicationConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Jwt      jwtConfig      `yaml:"jwt"`
}

func InitConfig(log *logrus.Logger) (*applicationConfig, error) {
	environment := viper.GetString("environment")
	if len(environment) == 0 {
		environment = "prod"
	}

	config := viper.New()
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType(configType)

	// для локальной разработки
	var applicationConfig applicationConfig
	switch environment {
	case "local":
		if err := config.ReadInConfig(); err != nil {
			return nil, err
		}
		if err := config.Unmarshal(&applicationConfig); err != nil {
			return nil, err
		}
	default:
		// для куба (все что не local)

		// сервис
		applicationConfig.Server.Port = os.Getenv("SERVICE_PORT")
		applicationConfig.Server.Host = os.Getenv("SERVICE_HOST")

		// бд
		applicationConfig.Database.Drivername = os.Getenv("DB_DRIVER_NAME")
		applicationConfig.Database.Host = os.Getenv("POSTGRES_HOST")
		applicationConfig.Database.Port = os.Getenv("POSTGRES_PORT")
		applicationConfig.Database.Username = os.Getenv("POSTGRES_USER")
		applicationConfig.Database.Password = os.Getenv("POSTGRES_PASSWORD")
		applicationConfig.Database.Dbname = os.Getenv("POSTGRES_DB")
		applicationConfig.Database.Sslmode = os.Getenv("DB_SSL_MODE")

		// jwtConfig
		applicationConfig.Jwt.AccessSecret = os.Getenv("JWT_ACCESS_SECRET")
		applicationConfig.Jwt.RefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	}

	if marshal, err := json.MarshalIndent(&applicationConfig, "", "  "); err == nil {
		log.Println(fmt.Sprintf("Load application configuration %s", marshal))
	}

	return &applicationConfig, nil
}
