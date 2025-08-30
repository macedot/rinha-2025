package config

import (
	"log"
	"rinha-2025-go/pkg/util"
	"strconv"
)

type Config struct {
	ServerSocket     string
	PollWorkerSize   int
	DatabaseURL      string
	DatabaseMaxConns int
}

func GetConfig() *Config {
	return &Config{
		ServerSocket: util.GetEnv("SOCKET_PATH"),
		DatabaseURL:  util.GetEnvOr("DATABASE_URL", "postgresql://rinha:rinha@postgres:5432/rinha?sslmode=disable"),
		PollWorkerSize: func() int {
			pollWorkerSize, err := strconv.Atoi(util.GetEnvOr("WORKERS_POOL_SIZE", "5"))
			if err != nil {
				log.Fatalf("Invalid value for WORKERS_POOL_SIZE: %v", err)
			}
			if pollWorkerSize < 1 {
				log.Fatalf("Invalid value for WORKERS_POOL_SIZE: %v", err)
			}
			return pollWorkerSize
		}(),
		DatabaseMaxConns: func() int {
			maxConns, err := strconv.Atoi(util.GetEnvOr("DB_MAX_CONNS", "25"))
			if err != nil {
				log.Fatalf("Invalid value for DB_MAX_CONNS: %v", err)
			}
			if maxConns < 1 {
				log.Fatalf("Invalid value for DB_MAX_CONNS: %v", err)
			}
			return maxConns
		}(),
	}
}
