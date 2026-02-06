/* Сервисы для управления записью процессов с высоким использованием ресурсов */
package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var (
	recordingActive  bool
	recordingSession *RecordingSession
	recordingMutex   sync.RWMutex
	recordingCancel  context.CancelFunc
)

/* Структура для хранения информации о сессии записи */
type RecordingSession struct {
	ID           int64
	CPUThreshold float64
	RAMThreshold float64
	Duration     int
	StartedAt    time.Time
	EndTime      time.Time
}

/*
Запускает новую сессию записи процессов с указанными порогами и длительностью

	Возвращает ID сессии или ошибку
*/
func StartRecording(cpuThreshold, ramThreshold float64, durationSec int) (int64, error) {
	db := GetDB()
	if db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	recordingMutex.Lock()
	defer recordingMutex.Unlock()

	if recordingActive {
		return 0, fmt.Errorf("запись уже активна")
	}

	startedAt := time.Now()
	endTime := startedAt.Add(time.Duration(durationSec) * time.Second)

	result, err := db.Exec(
		"INSERT INTO recording_sessions (started_at, ended_at, cpu_threshold, ram_threshold, duration_sec, status) VALUES (?, ?, ?, ?, ?, 'active')",
		startedAt.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"),
		cpuThreshold,
		ramThreshold,
		durationSec,
	)
	if err != nil {
		return 0, err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	recordingSession = &RecordingSession{
		ID:           sessionID,
		CPUThreshold: cpuThreshold,
		RAMThreshold: ramThreshold,
		Duration:     durationSec,
		StartedAt:    startedAt,
		EndTime:      endTime,
	}

	recordingActive = true

	ctx, cancel := context.WithCancel(context.Background())
	recordingCancel = cancel

	go recordProcessesLoop(ctx, db, sessionID, cpuThreshold, ramThreshold, endTime)

	return sessionID, nil
}

/* Останавливает активную сессию записи и обновляет статус в базе данных */
func StopRecording() error {
	recordingMutex.Lock()
	defer recordingMutex.Unlock()

	if !recordingActive {
		return fmt.Errorf("запись не активна")
	}

	if recordingCancel != nil {
		recordingCancel()
		recordingCancel = nil
	}

	db := GetDB()
	if db != nil && recordingSession != nil {
		db.Exec(
			"UPDATE recording_sessions SET status = 'stopped', ended_at = ? WHERE id = ?",
			time.Now().Format("2006-01-02 15:04:05"),
			recordingSession.ID,
		)
	}

	recordingActive = false
	recordingSession = nil

	return nil
}

/* Возвращает текущий статус записи и информацию о сессии если она активна */
func GetRecordingStatus() (bool, *RecordingSession) {
	recordingMutex.RLock()
	defer recordingMutex.RUnlock()
	return recordingActive, recordingSession
}

/*
Основной цикл записи процессов которые превышают пороги cpu и ram

	Записывает процессы в базу данных каждые 2 секунды до окончания времени сессии
*/
func recordProcessesLoop(ctx context.Context, db *sql.DB, sessionID int64, cpuThreshold, ramThreshold float64, endTime time.Time) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if time.Now().After(endTime) {
				StopRecording()
				return
			}

			procs, err := process.Processes()
			if err != nil {
				log.Printf("Ошибка получения процессов: %v", err)
				continue
			}

			for _, p := range procs {
				cpuPercent, err := p.CPUPercent()
				if err != nil {
					continue
				}

				memInfo, err := p.MemoryInfo()
				if err != nil {
					continue
				}

				memPercent, err := p.MemoryPercent()
				if err != nil {
					continue
				}

				if cpuPercent > cpuThreshold && float64(memPercent) > ramThreshold {
					saveRecordedProcess(db, sessionID, p, cpuPercent, float64(memPercent), memInfo.RSS)
				}
			}
		}
	}
}

/* Сохраняет информацию о процессе в базу данных если он превышает пороги */
func saveRecordedProcess(db *sql.DB, sessionID int64, p *process.Process, cpuPercent, memPercent float64, memRSS uint64) {
	name, _ := p.Name()
	exe, _ := p.Exe()
	cmdline, _ := p.Cmdline()
	username, _ := p.Username()

	_, err := db.Exec(
		"INSERT INTO recorded_processes (session_id, recorded_at, pid, name, cpu_percent, memory_percent, memory_rss, exe, cmdline, username) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		sessionID,
		time.Now().Format("2006-01-02 15:04:05"),
		p.Pid,
		name,
		cpuPercent,
		memPercent,
		memRSS,
		exe,
		cmdline,
		username,
	)
	if err != nil {
		log.Printf("Ошибка сохранения процесса: %v", err)
	}
}

/*
Получает список записанных процессов для указанной сессии с ограничением количества записей

	Возвращает массив словарей с информацией о процессах или ошибку
*/
func GetRecordedProcesses(sessionID int64, limit int) ([]map[string]interface{}, error) {
	db := GetDB()
	if db == nil {
		return nil, nil
	}

	query := "SELECT recorded_at, pid, name, cpu_percent, memory_percent, memory_rss, exe, cmdline, username FROM recorded_processes WHERE session_id = ? ORDER BY recorded_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
	}

	rows, err := db.Query(query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var timestamp string
		var pid int32
		var name, exe, cmdline, username string
		var cpuPercent, memPercent float64
		var memRSS int64

		if err := rows.Scan(&timestamp, &pid, &name, &cpuPercent, &memPercent, &memRSS, &exe, &cmdline, &username); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"timestamp":     timestamp,
			"pid":           pid,
			"name":          name,
			"cpuPercent":    cpuPercent,
			"memoryPercent": memPercent,
			"memoryRSS":     memRSS,
			"exe":           exe,
			"cmdline":       cmdline,
			"username":      username,
		})
	}

	return results, nil
}
