package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/bot"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/config"
	db "github.com/skewertoreversepolarity/telegramBotCrypto/internal/database"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/monitor"
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

	dbMonitor := monitor.New(db, telegramBot, cfg.WalletAddress, cfg.PolliInterval)

	go dbMonitor.Start()

	log.Println("Бот запущен. Ожидание событий...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Завершение работы бота...")

}
