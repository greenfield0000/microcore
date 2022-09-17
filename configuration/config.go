package configuration

import (
	"encoding/json"
	"errors"
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

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Drivername string `yaml:"drivername"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Dbname     string `yaml:"dbname"`
	Sslmode    string `yaml:"sslmode"`
}

type Jwt struct {
	Accesssecret  string `yaml:"accesssecret"`
	Refreshsecret string `yaml:"refreshsecret"`
}

type ApplicationConfig struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Jwt      Jwt      `yaml:"jwt"`
}

func InitConfig(log *logrus.Logger) (*ApplicationConfig, error) {
	environment := viper.GetString("environment")
	if len(environment) == 0 {
		return nil, errors.New("Flag environment no set!")
	}

	config := viper.New()
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType(configType)

	// для локальной разработки
	var applicationConfig ApplicationConfig
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

		// jwt
		applicationConfig.Jwt.Accesssecret = os.Getenv("JWT_ACCESS_SECRET")
		applicationConfig.Jwt.Refreshsecret = os.Getenv("JWT_REFRESH_SECRET")
	}

	if marshal, err := json.MarshalIndent(&applicationConfig, "", "  "); err == nil {
		log.Println(fmt.Sprintf("Load application configuration %s", marshal))
	}

	return &applicationConfig, nil
}
