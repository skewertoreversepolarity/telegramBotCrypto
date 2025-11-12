package bot

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/models"
	"gopkg.in/telebot.v3"
)

// Bot represents a Telegram bot.
type Bot struct {
	bot    *telebot.Bot
	chatID int64
}

func New(token string, chatID int64) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("Токен Telegram бота не указан")
	}

	if chatID == 0 {
		return nil, fmt.Errorf("Chat ID не указан")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})

	if err != nil {
		return nil, fmt.Errorf("Ошибка создания бота: %w", err)
	}

	info, err := bot.Raw("getMe", nil)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к Telegram API: %w", err)
	}
	log.Printf("Бот %s успешно создан\n", info)
	log.Printf("Информация о боте: %s", info)

	return &Bot{
		bot:    bot,
		chatID: chatID,
	}, nil

}

func (b *Bot) SendDepositNotification(deposit *models.Deposit) error {
	message := fmt.Sprintf(
		"💰 *Новый депозит*\n\n"+
			"🌐 Сеть: %s\n"+
			"💱 Валюта: %s\n"+
			"📤 От: `%s`\n"+
			"📥 К: `%s`\n"+
			"💵 Сумма: %.2f\n"+
			"🔗 Хеш: `%s`\n"+
			"📊 Блок: %d\n"+
			"📋 Результат: %s\n"+
			"⚡ Статус: %s\n"+
			"🕐 Время: %s",
		deposit.Network,
		deposit.Currency,
		deposit.FromAddress,
		deposit.ToAddress,
		deposit.Amount,
		deposit.TxnHash,
		deposit.BlockNumber,
		deposit.Result,
		deposit.Status,
		deposit.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	_, err := b.bot.Send(telebot.ChatID(b.chatID), message, telebot.ModeMarkdown)
	if err != nil {
		return fmt.Errorf("ошибка отправки уведомления о депозите: %w", err)
	}

	log.Printf("Отправлено уведомление о депозите ID: %d", deposit.ID)
	return nil
}

func (b *Bot) SendOutgoingNotification(outgoing *models.Outgoing) error {
	message := fmt.Sprintf(
		"💸 *Новый исходящий перевод*\n\n"+
			"🌐 Сеть: %s\n"+
			"💱 Валюта: %s\n"+
			"📤 От: `%s`\n"+
			"📥 К: `%s`\n"+
			"💵 Сумма: %.2f\n"+
			"💸 Комиссия: %.2f\n"+
			"💰 Итого: %.2f\n"+
			"🔗 Хеш: `%s`\n"+
			"⚡ Статус: %s\n"+
			"🕐 Время: %s",
		outgoing.Network,
		outgoing.Currency,
		outgoing.From,
		outgoing.To,
		outgoing.Amount,
		outgoing.Commission,
		outgoing.TotalAmount,
		getStringValue(outgoing.TranHash),
		outgoing.Status,
		outgoing.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	_, err := b.bot.Send(telebot.ChatID(b.chatID), message, telebot.ModeMarkdown)
	if err != nil {
		return fmt.Errorf("ошибка отправки уведомления об исходящем переводе: %w", err)
	}

	log.Printf("Отправлено уведомление об исходящем переводе ID: %d", outgoing.ID)
	return nil
}

func (b *Bot) SendBalanceNotification(message string) error {
	_, err := b.bot.Send(telebot.ChatID(b.chatID), message, telebot.ModeMarkdown)
	if err != nil {
		return fmt.Errorf("ошибка отправки уведомления о балансе: %w", err)
	}

	log.Printf("Отправлено уведомление о балансе")
	return nil
}

func getStringValue(s *string) string {
	if s == nil {
		return "N/A"
	}
	return *s
}
