package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrEmptyEnv = errors.New("empty env variable, please provide all tokens")
)

type Config interface {
	CsgoMarketToken() string
	LisSkinsToken() string
	TelegramToken() string
}

type envConfig struct {
	csgoMarketToken string
	lisSkinsToken   string
	telegramToken   string
}

func (e envConfig) TelegramToken() string   { return e.telegramToken }
func (e envConfig) CsgoMarketToken() string { return e.csgoMarketToken }
func (e envConfig) LisSkinsToken() string   { return e.lisSkinsToken }

func New() (Config, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("cant load config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid .env: %w", err)
	}

	return cfg, nil
}

func loadConfig() (*envConfig, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, fmt.Errorf("cant load .env: %w", err)
	}

	return &envConfig{
		csgoMarketToken: os.Getenv("CSMARKET"),
		lisSkinsToken:   os.Getenv("LISSKINS"),
		telegramToken:   os.Getenv("TELEGRAM_BOT"),
	}, nil
}

func (e *envConfig) validate() error {
	if e.csgoMarketToken == "" || e.lisSkinsToken == "" || e.telegramToken == "" {
		return ErrEmptyEnv
	}

	return nil
}
