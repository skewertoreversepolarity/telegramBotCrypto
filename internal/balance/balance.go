package balance

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

const USDTContractTRC20 = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"

func GetUSDTBalsnce(address string) (balanceUSDT float64, rawUnits *big.Int, err error) {
	addr := strings.TrimSpace(address)
	if addr == "" {
		return 0, nil, nil
	}
	if !strings.HasPrefix(addr, "T") || len(addr) < 25 {
		WriteLog("Некорректный адрес TRON: " + addr)
	}

	url := fmt.Sprintf("https://apilist.tronscan.org/api/account?address=%s", addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Ошибка создания HTTP запроса:", err)
		return 0, nil, err
	}
	req.Header.Set("User-Agent", "crypto-bot/1.0 (+balance-checker)")
	client := &http.Client{Timeout: 12 * time.Second}

	WriteLog("TronScan Get: " + url)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка выполнения HTTP запроса:", err)
		return 0, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка ответа от TronScan: %s", resp.Status)
		return 0, nil, fmt.Errorf("Ошибка ответа от TronScan: %s", resp.Status)
	}
	var js map[string]any
	if err = json.NewDecoder(resp.Body).Decode(&js); err != nil {
		log.Println("Ошибка декодирования JSON ответа:", err)
		return 0, nil, err
	}

	// A) js["trc20"]
	if arr, ok := js["trc20"].([]any); ok {
		if bal, dec, found := extractUSDTFromTrc20Array(arr); found {
			return humanizeTRC20(bal, dec)
		}
	}
	// B) js["tokens"]
	if arr, ok := js["tokens"].([]any); ok {
		if bal, dec, found := extractUSDTFromTokensArray(arr); found {
			return humanizeTRC20(bal, dec)
		}
	}
	// C) js["trc20token_balances"]
	if arr, ok := js["trc20token_balances"].([]any); ok {
		if bal, dec, found := extractUSDTFromTokensArray(arr); found {
			return humanizeTRC20(bal, dec)
		}
	}

	// не нашли — считаем 0
	WriteLog(fmt.Sprintf("ℹ️ USDT не найден у адреса %s, баланс = 0", addr))
	return 0, big.NewInt(0), nil
}

func WriteLog(message string) {
	logFile := "crypto_bot.log"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Ошибка открытия файла лога:", err)
		return
	}
	defer file.Close()

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s - %s\n", timeStamp, message)
	file.WriteString(logEntry)
}

// --- Вспомогательные парсеры ---

// extractUSDTFromTrc20Array ищет USDT в массиве js["trc20"].
// Возвращает (rawBalance, decimals, found)
func extractUSDTFromTrc20Array(arr []any) (*big.Int, int, bool) {
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		// вариант {"TR7NH...":"123456"}
		if len(m) == 1 {
			for k, v := range m {
				if strings.EqualFold(k, USDTContractTRC20) {
					if s, ok := v.(string); ok {
						if bi, ok := new(big.Int).SetString(s, 10); ok {
							return bi, 6, true
						}
					}
				}
			}
			continue
		}
		// вариант с полями
		contract := str(m["tokenId"])
		abbr := strings.ToUpper(str(m["tokenAbbr"]))
		if strings.EqualFold(contract, USDTContractTRC20) || abbr == "USDT" {
			balStr := str(m["balance"])
			dec := asIntDefault(m["tokenDecimal"], 6)
			if bi, ok := new(big.Int).SetString(balStr, 10); ok {
				return bi, dec, true
			}
		}
	}
	return nil, 0, false
}

// extractUSDTFromTokensArray ищет USDT в массиве js["tokens"] или js["trc20token_balances"].
// Возвращает (rawBalance, decimals, found)
func extractUSDTFromTokensArray(arr []any) (*big.Int, int, bool) {
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		contract := str(m["tokenId"])
		abbr := strings.ToUpper(str(m["tokenAbbr"]))
		name := strings.ToUpper(str(m["tokenName"]))
		if strings.EqualFold(contract, USDTContractTRC20) || abbr == "USDT" || strings.Contains(name, "TETHER") {
			balStr := str(m["balance"])
			dec := asIntDefault(m["tokenDecimal"], 6)
			if bi, ok := new(big.Int).SetString(balStr, 10); ok {
				return bi, dec, true
			}
		}
	}
	return nil, 0, false
}

func str(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	case json.Number:
		return t.String()
	default:
		return fmt.Sprint(v)
	}
}

func asIntDefault(v any, def int) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case string:
		if t == "" {
			return def
		}
		var x int
		if _, err := fmt.Sscanf(t, "%d", &x); err == nil {
			return x
		}
	}
	return def
}

// humanizeTRC20 переводит из минимальных единиц в человекочитаемый float64 с учётом decimals
func humanizeTRC20(raw *big.Int, decimals int) (float64, *big.Int, error) {
	if raw == nil {
		raw = big.NewInt(0)
	}
	if decimals < 0 || decimals > 36 {
		decimals = 6
	}
	den := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	num := new(big.Float).SetInt(raw)
	f, _ := new(big.Float).Quo(num, den).Float64()

	pow := math.Pow10(decimals)
	f = math.Round(f*pow) / pow

	return f, raw, nil
}

// PrintUSDTBalance печатает в stdout и пишет в ваши логи
func PrintUSDTBalance(address string) {
	bal, raw, err := GetUSDTBalsnce(address)
	if err != nil {
		log.Printf("❌ USDT balance error for %s: %v\n", address, err)
		WriteLog("❌ USDT balance error: " + err.Error())
		return
	}
	msg := fmt.Sprintf("💰 USDT balance for %s: %.6f (raw=%s)", address, bal, raw.String())
	fmt.Println(msg) // консоль
	WriteLog(msg)    // файл логов
}
