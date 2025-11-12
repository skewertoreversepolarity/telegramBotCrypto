package main

import (
	"log"

	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/bot"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/config"
	db "github.com/skewertoreversepolarity/telegramBotCrypto/internal/database"
)

// echo "# telegramBotCrypto" >> README.md
// git init
// git add README.md
// git commit -m "first commit"
// git branch -M main
// git remote add origin git@github.com:skewertoreversepolarity/telegramBotCrypto.git
// git push -u origin main
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Загрузка конфигураций: %+v\n", cfg)
	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()
	log.Println("Успешное подключение к БД")

	telegramBot, err := bot.New(cfg.TelegramToken, cfg.ChatID)
	if err != nil {
		log.Fatal("Ошибка создания Telegram бота:", err)
	}

}
