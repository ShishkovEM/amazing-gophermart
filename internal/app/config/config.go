package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	Database      string `env:"DATABASE_URI"`
	AccrualSystem string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (cfg *Config) Parse() {

	// Считывает конфигурацию из переменных окружения
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	// Если передана конфигурация флагами - используем её
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.StringVar(&cfg.AccrualSystem, "r", cfg.AccrualSystem, "address of the accrual system")
	flag.Parse()

	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
}
