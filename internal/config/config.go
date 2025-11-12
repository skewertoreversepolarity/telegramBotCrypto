package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	TelegramToken string
	ChatID        int64
	PolliInterval int
	WalletAddress string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	chatID, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	polliInterval, err := strconv.Atoi(os.Getenv("POLLI_INTERVAL"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		ChatID:        chatID,
		PolliInterval: polliInterval,
		WalletAddress: os.Getenv("WALLET_ADDRESS"),
	}, nil
}
