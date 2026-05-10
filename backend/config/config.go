package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	configFile := fmt.Sprintf("config.%s.env", env)
	viper.SetConfigName(configFile)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Set defaults
	// viper.SetDefault("DB_HOST", "localhost")
	// viper.SetDefault("DB_PORT", "5432")
	// viper.SetDefault("DB_USER", "postgres")
	// viper.SetDefault("DB_PASSWORD", "postgres")
	// viper.SetDefault("DB_NAME", "movie_reck")

	// Read from environment variables if present
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: error reading config file, relying on environment variables: %v", err)
	} else {
		log.Printf("Successfully loaded config from %s", viper.ConfigFileUsed())
	}
}
