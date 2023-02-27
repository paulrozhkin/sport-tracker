package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
}

type ServerConfigurations struct {
	Port int
}

type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfigurations() (*Configurations, error) {
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

	configuration := &Configurations{}

	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return configuration, nil
}
