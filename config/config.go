package config

import (
	"fmt"
	"os"
)

// Configuration holds various configuration settings
type Configuration struct {
	Env           Env42
	PostgreeDBDsn string
	RedisHost     string
	RedisPassword string
}

// GetConfig returns the current configuration
func GetConfig() Configuration {
	return Configuration{
		Env:           CurrentEnv(),
		PostgreeDBDsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
