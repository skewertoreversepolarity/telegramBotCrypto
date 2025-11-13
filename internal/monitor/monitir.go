package monitor

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/balance"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/bot"
	"github.com/skewertoreversepolarity/telegramBotCrypto/internal/models"
)

type Monitor struct {
	db            *sql.DB
	bot           *bot.Bot
	walletAddress string
	pollInterval  int
}

func New(db *sql.DB, bot *bot.Bot, walletAddress string, pollInterval int) *Monitor {
	if pollInterval <= 0 {
		pollInterval = 5
	}
	return &Monitor{
		db:            db,
		bot:           bot,
		walletAddress: walletAddress,
		pollInterval:  pollInterval,
	}
}

func (m *Monitor) Start() {
	log.Println("Мониторинг запущен")

}

func (m *Monitor) getLastDepositID() (int, error) {
	var id int
	err := m.db.QueryRow("SELECT COALESCE(MAX(id),0)FROM deposit").Scan(&id)
	return id, err
}

func (m *Monitor) getLastOutgoingID() (int, error) {
	var id int
	err := m.db.QueryRow("SELECT COALESCE(MAX(id),0)FROM outgoing").Scan(&id)
	return id, err
}

func (m *Monitor) getNewDeposits(lastID int) ([]models.Deposit, error) {
	query := `
		SELECT id, network, currency, from_address, to_address, amount,
		 txn_hash, block_number, result, status, created_at, updated_at, tran_code
		 FROM deposit
		 WHERE id > $1
		 ORDER BY id ASC
	`
	rows, err := m.db.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deposits []models.Deposit
	for rows.Next() {
		var deposit models.Deposit
		err := rows.Scan(
			&deposit.ID,
			&deposit.Network,
			&deposit.Currency,
			&deposit.FromAddress,
			&deposit.ToAddress,
			&deposit.Amount,
			&deposit.TxnHash,
			&deposit.BlockNumber,
			&deposit.Result,
			&deposit.Status,
			&deposit.CreatedAt,
			&deposit.UpdatedAt,
			&deposit.TranCode,
		)
		if err != nil {
			return nil, err
		}

		deposits = append(deposits, deposit)

	}
	return deposits, nil
}

func (m *Monitor) getNewOutgoings(lastID int) ([]models.Outgoing, error) {
	query := `
		SELECT id, "from", "to", currency, network, amount, commission,
		       trancode, tran_hash, status, response, created_at, updated_at,
		       total_amount, finoper
		FROM outgoing 
		WHERE id > $1 
		ORDER BY id ASC
	`
	rows, err := m.db.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var outgoings []models.Outgoing
	for rows.Next() {
		var outgoing models.Outgoing
		err := rows.Scan(
			&outgoing.ID, &outgoing.From, &outgoing.To, &outgoing.Currency,
			&outgoing.Network, &outgoing.Amount, &outgoing.Commission,
			&outgoing.TranCode, &outgoing.TranHash, &outgoing.Status,
			&outgoing.Response, &outgoing.CreatedAt, &outgoing.UpdatedAt,
			&outgoing.TotalAmount, &outgoing.Finoper,
		)
		if err != nil {
			return nil, err
		}
		outgoings = append(outgoings, outgoing)
	}
	return outgoings, nil
}

func (m *Monitor) checkAndSendBalance() {
	if m.walletAddress == "" {
		log.Println("Адрес кошелька не задан, пропуск проверки баланса")
		return
	}
	log.Println("Проверка баланса кошелька:", m.walletAddress)

	usdtBalance, rawUnits, err := balance.GetUSDTBalsnce(m.walletAddress)
	if err != nil {
		log.Println("Ошибка получения баланса USDT:", err)
		balance.WriteLog(fmt.Sprintf("❌ Ошибка получения баланса USDT для %s: %v", m.walletAddress, err))
		return
	}
	// Отправляем уведомление о балансе
	message := fmt.Sprintf(
		"💰 *Текущий баланс USDT*\n\n"+
			"🏦 Адрес: `%s`\n"+
			"💵 Баланс: %.6f USDT\n"+
			"🔢 Сырые единицы: %s\n"+
			"🕐 Время: %s",
		m.walletAddress,
		usdtBalance,
		rawUnits.String(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	err = m.bot.SendBalanceNotification(message)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления о балансе: %v", err)
		balance.WriteLog(fmt.Sprintf("❌ Ошибка отправки уведомления о балансе: %v", err))
	} else {
		log.Printf("✅ Уведомление о балансе отправлено успешно: %.6f USDT", usdtBalance)
		balance.WriteLog(fmt.Sprintf("✅ Уведомление о балансе отправлено: %.6f USDT", usdtBalance))
	}

}
