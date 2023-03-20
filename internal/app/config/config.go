package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	Database                string `env:"DATABASE_URI"`
	AccrualSystem           string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey               string `env:"SECRET_KEY" envDefault:"G0pher"`
	RequestTimeout          string `env:"REQUEST_TIMEOUT_SECONDS" envDefault:"60s"`
	ExpBackOffInitialAmount string `env:"EXP_BACKOFF_MILLIS" envDefault:"100ms"`
	CoolDownDuration        string `env:"COOLDOWN_DURATION_SECONDS" envDefault:"60s"`
	ServerReadTimeout       string `env:"SERVER_READ_TIMEOUT" envDefault:"60s"`
	ServerWriteTimeout      string `env:"SERVER_WRITE_TIMEOUT" envDefault:"60s"`
	CookieLifeTime          string `env:"COOKIE_LIFETIME" envDefault:"8760h"`
}

func (cfg *Config) Parse() {

	// Считывает конфигурацию из переменных окружения
	err := env.Parse(cfg)

	if err != nil {
		log.Fatal(err)
	}

	// Если передана конфигурация флагами - используем её
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.StringVar(&cfg.AccrualSystem, "r", cfg.AccrualSystem, "address of the accrual system")
	flag.StringVar(&cfg.SecretKey, "s", cfg.SecretKey, "secret key")
	flag.StringVar(&cfg.RequestTimeout, "t", cfg.RequestTimeout, "60s")
	flag.StringVar(&cfg.ExpBackOffInitialAmount, "e", cfg.ExpBackOffInitialAmount, "100ms")
	flag.StringVar(&cfg.CoolDownDuration, "c", cfg.CoolDownDuration, "60s")
	flag.StringVar(&cfg.ServerReadTimeout, "o", cfg.ServerReadTimeout, "60s")
	flag.StringVar(&cfg.ServerWriteTimeout, "k", cfg.ServerWriteTimeout, "60s")
	flag.StringVar(&cfg.CookieLifeTime, "l", cfg.CookieLifeTime, "8760h")
	flag.Parse()

	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
}
