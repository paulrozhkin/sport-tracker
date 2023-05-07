package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const DefaultRequestTimeoutSeconds = 60

type Configuration struct {
	Server        ServerConfigurations
	Database      DatabaseConfigurations
	JwtSigningKey string `yaml:"jwtSigningKey"`
}

type ServerConfigurations struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	RequestTimeoutSeconds int    `yaml:"requestTimeoutSeconds"`
	DisableCORS           bool   `yaml:"disableCORS"`
}

type DatabaseConfigurations struct {
	DBName       string `yaml:"dbname"`
	DBUser       string `yaml:"dbuser"`
	DBPassword   string `yaml:"dbpassword"`
	DBConnection string `yaml:"dbconnection"`
	DBSslMode    string `yaml:"dbsslmode"`
	DataSeeding  bool   `yaml:"dbseeding"`
}

func LoadConfigurations() (*Configuration, error) {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TRACKER")

	_ = viper.BindEnv("server.port", "TRACKER_SERVER_PORT")
	_ = viper.BindEnv("server.host", "TRACKER_SERVER_HOST")
	_ = viper.BindEnv("server.requestTimeoutSeconds", "TRACKER_SERVER_TIMEOUT")
	_ = viper.BindEnv("server.disableCORS", "TRACKER_SERVER_DISABLE_CORS")
	_ = viper.BindEnv("database.dbname", "TRACKER_DATABASE_NAME")
	_ = viper.BindEnv("database.dbuser", "TRACKER_DATABASE_USER")
	_ = viper.BindEnv("database.dbpassword", "TRACKER_DATABASE_PASSWORD")
	_ = viper.BindEnv("database.dbconnection", "TRACKER_DATABASE_CONNECTION")
	_ = viper.BindEnv("database.dbsslmode", "TRACKER_DATABASE_SSL")
	_ = viper.BindEnv("database.dbseeding", "TRACKER_DATABASE_SEEDING")

	configuration := &Configuration{
		Server: ServerConfigurations{
			RequestTimeoutSeconds: DefaultRequestTimeoutSeconds,
		},
		Database: DatabaseConfigurations{
			DataSeeding: false,
			DBSslMode:   "disable",
		},
	}

	var err error
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	confString, _ := json.MarshalIndent(configuration, "", " ")
	zap.S().Info("Configuration:\n", string(confString))
	validationError := configuration.validate()
	if validationError != nil {
		return nil, validationError
	}
	return configuration, nil
}

func (c *Configuration) validate() error {
	if c.JwtSigningKey == "" {
		return errors.New("JwtSigningKey null or empty")
	}
	return nil
}
