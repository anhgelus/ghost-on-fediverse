package src

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

var (
	Client *redis.Client
	ctx    = context.Background()

	ErrEnvNotSet = errors.New("env var not set")
)

const (
	EnvRedisHost     = "FED_REDIS_HOST"
	EnvRedisPort     = "FED_REDIS_PORT"
	EnvRedisPassword = "FED_REDIS_PASSWORD"
)

func ConnectRedis() error {
	host, b := os.LookupEnv(EnvRedisHost)
	if !b {
		return ErrEnvNotSet
	}

	portStr, b := os.LookupEnv(EnvRedisPort)
	if !b {
		return ErrEnvNotSet
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	password := os.Getenv(EnvRedisPassword)
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return Client.Ping(ctx).Err()
}
