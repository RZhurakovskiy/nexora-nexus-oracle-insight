/* Сервисы для работы с алертами системы */
package services

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

/* Структура для хранения информации об алерте */
type Alert struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Type         string    `json:"type"`         // "cpu" или "memory"
	Threshold    float64   `json:"threshold"`    // Пороговое значение
	CurrentValue float64   `json:"currentValue"` // Текущее значение
	Message      string    `json:"message"`
	Acknowledged bool      `json:"acknowledged"`
}

// Сохраняет новый алерт в базу данных с указанным типом, порогом, текущим значением и сообщением
func SaveAlert(alertType string, threshold, currentValue float64, message string) error {
	var db *sql.DB = GetDB()
	if db == nil {
		return nil
	}

	_, err := db.Exec(
		"INSERT INTO alerts (created_at, type, threshold, current_value, message, acknowledged) VALUES (?, ?, ?, ?, ?, 0)",
		time.Now().Format("2006-01-02 15:04:05"),
		alertType,
		threshold,
		currentValue,
		message,
	)
	return err
}

/*
Получает список алертов с возможностью фильтрации по неподтвержденным и ограничения количества

	Возвращает массив алертов или ошибку
*/
func GetAlerts(limit int, unacknowledgedOnly bool) ([]Alert, error) {
	db := GetDB()
	if db == nil {
		return nil, nil
	}

	query := "SELECT id, created_at, type, threshold, current_value, message, acknowledged FROM alerts"
	if unacknowledgedOnly {
		query += " WHERE acknowledged = 0"
	}
	query += " ORDER BY created_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
	}

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		var createdAtStr string
		var acknowledged int

		if err := rows.Scan(&a.ID, &createdAtStr, &a.Type, &a.Threshold, &a.CurrentValue, &a.Message, &acknowledged); err != nil {
			return nil, err
		}

		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			a.CreatedAt = parsed
		}
		a.Acknowledged = acknowledged == 1

		alerts = append(alerts, a)
	}

	return alerts, nil
}

/* Отмечает алерт как подтвержденный по его ID в базе данных */
func AcknowledgeAlert(id int64) error {
	db := GetDB()
	if db == nil {
		return nil
	}

	_, err := db.Exec("UPDATE alerts SET acknowledged = 1 WHERE id = ?", id)
	return err
}

var (
	alertCPUThreshold    float64
	alertMemoryThreshold float64
	alertMutex           sync.RWMutex
)

func init() {
	alertCPUThreshold = 0
	alertMemoryThreshold = 0
}

/* Устанавливает пороги для CPU и памяти, при превышении которых будут создаваться алерты */
func SetAlertThresholds(cpuThreshold, memoryThreshold float64) {
	alertMutex.Lock()
	defer alertMutex.Unlock()
	alertCPUThreshold = cpuThreshold
	alertMemoryThreshold = memoryThreshold
}

/* Возвращает текущие пороги для CPU и памяти в виде словаря */
func GetAlertThresholds() map[string]float64 {
	alertMutex.RLock()
	defer alertMutex.RUnlock()
	return map[string]float64{
		"cpuThreshold":    alertCPUThreshold,
		"memoryThreshold": alertMemoryThreshold,
	}
}

/* Проверяет текущие значения CPU и памяти и создает алерты если они превышают установленные пороги */
func CheckAlerts(currentCPU, currentMemory float64) {
	alertMutex.RLock()
	cpuThreshold := alertCPUThreshold
	memoryThreshold := alertMemoryThreshold
	alertMutex.RUnlock()

	if cpuThreshold > 0 && currentCPU > cpuThreshold {
		message := fmt.Sprintf("Превышен порог CPU: %.2f%% (порог: %.2f%%)", currentCPU, cpuThreshold)
		SaveAlert("cpu", cpuThreshold, currentCPU, message)
	}

	if memoryThreshold > 0 && currentMemory > memoryThreshold {
		message := fmt.Sprintf("Превышен порог памяти: %.2f%% (порог: %.2f%%)", currentMemory, memoryThreshold)
		SaveAlert("memory", memoryThreshold, currentMemory, message)
	}
}
