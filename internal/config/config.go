package config

import (
	"fmt"

	"github.com/mistic0xb/smolurl/internal/logger"
	"github.com/mistic0xb/smolurl/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Database DatabaseConfig `mapstructure:"database" validate:"required"`
	Redis    RedisConfig    `mapstructure:"redis" validate:"required"`
}

type ServerConfig struct {
	Port               string   `mapstructure:"port" validate:"required"`
	ReadTimeout        int      `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout       int      `mapstructure:"write_timeout" validate:"required"`
	IdleTimeout        int      `mapstructure:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `mapstructure:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required"`
	User            string `mapstructure:"user" validate:"required"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name" validate:"required"`
	SSLMode         string `mapstructure:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" validate:"required"`
}

type RedisConfig struct {
	Address string `mapstructure:"address" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := logger.NewLogger()

	// Set config file path
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("SMOLURL")
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Error().Err(err).Msg("No .env file found, using env vars only")
		return nil, err
	}

	// Unmarshal config into struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
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
