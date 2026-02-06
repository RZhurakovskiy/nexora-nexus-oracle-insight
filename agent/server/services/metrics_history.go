// Сервисы для работы с историей метрик в базе данных
package services

import (
	"database/sql"
	"sync"
	"time"
)

var (
	dbInstance *sql.DB
	dbMutex    sync.RWMutex
)

/* Устанавливает экземпляр базы данных для использования во всех сервисах */
func SetDB(db *sql.DB) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	dbInstance = db
}

func GetDB() *sql.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return dbInstance
}

/* Сохраняет метрики CPU и памяти в базу данных с текущей временной меткой */
func SaveMetricsHistory(cpuPercent, memoryPercent float64, memoryUsedMB, memoryTotalMB uint64) error {
	db := GetDB()
	if db == nil {
		return nil
	}
	_, err := db.Exec(
		"INSERT INTO metrics_history (timestamp, cpu_percent, memory_percent, memory_used_mb, memory_total_mb) VALUES (?, ?, ?, ?, ?)",
		time.Now().Format("2006-01-02 15:04:05"),
		cpuPercent,
		memoryPercent,
		memoryUsedMB,
		memoryTotalMB,
	)
	return err
}

/*
Получает историю метрик за указанный период времени с возможностью ограничения количества записей

	Возвращает массив словарей с метриками или ошибку
*/
func GetMetricsHistory(from, to time.Time, limit int) ([]map[string]interface{}, error) {
	db := GetDB()
	if db == nil {
		return nil, nil
	}
	var rows *sql.Rows
	var err error

	if limit > 0 {
		query := `
			SELECT timestamp, cpu_percent, memory_percent, memory_used_mb, memory_total_mb
			FROM metrics_history
			WHERE timestamp >= ? AND timestamp <= ?
			ORDER BY timestamp DESC
			LIMIT ?
		`
		rows, err = db.Query(query, from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"), limit)
	} else {
		query := `
			SELECT timestamp, cpu_percent, memory_percent, memory_used_mb, memory_total_mb
			FROM metrics_history
			WHERE timestamp >= ? AND timestamp <= ?
			ORDER BY timestamp DESC
		`
		rows, err = db.Query(query, from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var timestamp string
		var cpuPercent, memoryPercent float64
		var memoryUsedMB, memoryTotalMB int64

		if err := rows.Scan(&timestamp, &cpuPercent, &memoryPercent, &memoryUsedMB, &memoryTotalMB); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"timestamp":     timestamp,
			"cpuPercent":    cpuPercent,
			"memoryPercent": memoryPercent,
			"memoryUsedMB":  memoryUsedMB,
			"memoryTotalMB": memoryTotalMB,
		})
	}

	return results, nil
}

/* Удаляет старые записи метрик старше указанного количества дней */
func CleanOldMetricsHistory(daysToKeep int) error {
	db := GetDB()
	if db == nil {
		return nil
	}
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	_, err := db.Exec("DELETE FROM metrics_history WHERE timestamp < ?", cutoffDate.Format("2006-01-02 15:04:05"))
	return err
}

/* Удаляет все записи метрик из базы данных */
func ClearMetricsHistory() error {
	db := GetDB()
	if db == nil {
		return nil
	}
	_, err := db.Exec("DELETE FROM metrics_history")
	return err
}
