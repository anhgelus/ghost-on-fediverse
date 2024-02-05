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
	EnvHost     = "FED_HOST"
	EnvPort     = "FED_PORT"
	EnvPassword = "FED_PASSWORD"
)

func Connect() error {
	host, b := os.LookupEnv(EnvHost)
	if !b {
		return ErrEnvNotSet
	}

	portStr, b := os.LookupEnv(EnvPort)
	if !b {
		return ErrEnvNotSet
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	password := os.Getenv(EnvPassword)
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return Client.Ping(ctx).Err()
}
