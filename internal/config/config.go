package config

import (
	"errors"
	"fmt"
	"log"
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
	CsfloatCookie() string
}

type envConfig struct {
	csgoMarketToken string
	lisSkinsToken   string
	telegramToken   string
	csfloatCookie   string
}

func (e envConfig) TelegramToken() string   { return e.telegramToken }
func (e envConfig) CsgoMarketToken() string { return e.csgoMarketToken }
func (e envConfig) LisSkinsToken() string   { return e.lisSkinsToken }
func (e envConfig) CsfloatCookie() string   { return e.csfloatCookie }

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
	if os.Getenv("IS_DOCKER") == "" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &envConfig{
		csgoMarketToken: os.Getenv("CSMARKET"),
		lisSkinsToken:   os.Getenv("LISSKINS"),
		telegramToken:   os.Getenv("TELEGRAM_BOT"),
		csfloatCookie:   os.Getenv("CSFLOAT_COOKIE"),
	}, nil
}

func (e *envConfig) validate() error {
	if e.csgoMarketToken == "" || e.lisSkinsToken == "" || e.telegramToken == "" || e.csfloatCookie == "" {
		return ErrEmptyEnv
	}

	return nil
}
