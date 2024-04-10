package db

import (
	"fmt"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreeDatabase(dsn string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error initializing database:%v", err)
	}

	return db, nil
}

func InitRedis(host string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: host,
		// Password: password, // TODO for some reason the auth is failing with pass
		// DB:       16,
	})

	// Ping Redis to check if the connection is working
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Error initializing redis:%v", err)
	}

	return client, nil
}
