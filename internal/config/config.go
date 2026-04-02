package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mistic0xb/smolurl/internal/utils"
	"github.com/rs/zerolog"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig   `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Redis    RedisConfig    `koanf:"redis" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address  string `koanf:"address" validate:"required"`
	Password string `koanf:"password" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	var k = koanf.New(".")

	// load .env file into OS env vars 
	if err := godotenv.Load(); err != nil {
		logger.Warn().Msg(".env file not found, using OS env vars")
	}

	// load actual env vars (overrides .env)
	k.Load(env.Provider("", ".", func(s string) string {
		s = strings.ToLower(s)
		parts := strings.SplitN(s, "_", 2)
		if len(parts) == 2 {
			return parts[0] + "." + parts[1]
		}
		return s
	}), nil)

	// Unmarshal config into struct
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		logger.Error().Err(err).Msg("Unable to unmarshal to config struct")
		return nil, err
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid configuration")
	}

	return cfg, nil
}

// TODO: validate env
func (c *Config) Validate() error { return nil }

// Print displays the config (for debugging)
func (c *Config) Print() {
	fmt.Println("=== SmolURL Configuration ===")

	utils.PrintJSON("config:", c)

	fmt.Println("===================================")
}
