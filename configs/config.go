package config

import "github.com/spf13/viper"

// server config struct
type ServerConfig struct {
	Port string `mapstructure:"port"`
}


// database config struct
type DatabaseConfig struct {
	URI string `mapstructure:"uri"`
	DBName string `mapstructure:"dbname"`
}


// config struct 
type Config struct {
	Server  ServerConfig
	Database DatabaseConfig
}


// Load config from Config.yaml file
func LoadConfig() (*Config , error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config 

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config , nil
}