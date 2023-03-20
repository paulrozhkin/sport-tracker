package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
}

type ServerConfigurations struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

	configuration := &Configuration{}

	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	confString, _ := json.MarshalIndent(configuration, "", " ")
	zap.S().Info("Configuration:\n", string(confString))
	return configuration, nil
}
