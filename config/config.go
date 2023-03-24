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
}

type DatabaseConfigurations struct {
	DBName       string `yaml:"dbname"`
	DBUser       string `yaml:"dbuser"`
	DBPassword   string `yaml:"dbpassword"`
	DBConnection string `yaml:"dbconnection"`
	DBSslMode    string `yaml:"dbsslmode"`
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

	err := viper.BindEnv("server.port", "TRACKER_SERVER_PORT")
	if err != nil {
		return nil, fmt.Errorf("failed bind env due to: %v", err)
	}

	configuration := &Configuration{
		Server: ServerConfigurations{
			RequestTimeoutSeconds: DefaultRequestTimeoutSeconds,
		},
	}

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
