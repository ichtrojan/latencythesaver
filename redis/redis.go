package database

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func ConnectRedis(host string, port string, pass string, scheme string) error {
	if pass == "null" {
		pass = ""
	}

	var client *redis.Client

	if scheme == "tls" {
		client = redis.NewClient(&redis.Options{
			TLSConfig:  &tls.Config{},
			Addr:       fmt.Sprintf("%s:%s", host, port),
			Password:   pass,
			DB:         0,
			MaxRetries: 3,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: pass,
			DB:       0,
		})
	}

	_, err := client.Ping().Result()

	if err != nil {
		return err
	}

	Redis = client

	return nil
}
