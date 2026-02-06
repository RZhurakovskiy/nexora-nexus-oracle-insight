package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/RZhurakovskiy/agent/server/services"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v4/net"
)

/* Структура для хранения метрик CPU в кэше */
type cpuPayload struct {
	CPU       float64 `json:"cpu"`
	Timestamp string  `json:"timestamp"`
}

/* Структура для хранения метрик памяти в кэше */
type memoryPayload struct {
	MemoryUsage float64 `json:"memoryUsage"`
	UsedMB      uint64  `json:"usedMB"`
	TotalMemory uint64  `json:"totalMemory"`
	Timestamp   string  `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	cpuCache          cpuPayload
	memCache          memoryPayload
	procsCache        []models.ProcessInfo
	cacheMutex        sync.RWMutex
	cacheCtx          context.Context
	cacheCancel       context.CancelFunc
	monitoringEnabled bool
	monitoringMutex   sync.RWMutex
)

func init() {
	monitoringEnabled = false
}

/* Основной цикл обновления кэша метрик с разными интервалами для cpu, памяти и процессов */
/* Останавливается при получении сигнала отмены через контекст */
func updateCacheLoop(ctx context.Context) {

	cpuTicker := time.NewTicker(1 * time.Second)
	defer cpuTicker.Stop()

	memTicker := time.NewTicker(3 * time.Second)
	defer memTicker.Stop()

	procTicker := time.NewTicker(5 * time.Second)
	defer procTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cpuTicker.C:
			monitoringMutex.RLock()
			enabled := monitoringEnabled
			monitoringMutex.RUnlock()
			if enabled {
				updateCPUMetrics()
			}
		case <-memTicker.C:
			monitoringMutex.RLock()
			enabled := monitoringEnabled
			monitoringMutex.RUnlock()
			if enabled {
				updateMemoryMetrics()
			}
		case <-procTicker.C:
			monitoringMutex.RLock()
			enabled := monitoringEnabled
			monitoringMutex.RUnlock()
			if enabled {
				updateProcessMetrics()
			}
		}
	}
}

/* Обновляет кэш метрик cpu каждую секунду */
func updateCPUMetrics() {
	if usage, err := getmetrics.UsageCPU(100 * time.Millisecond); err == nil {
		cacheMutex.Lock()
		cpuCache = cpuPayload{
			CPU:       usage,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		cacheMutex.Unlock()
	} else {
		log.Printf("Ошибка обновления кэша CPU: %v", err)
	}
}

/* Обновляет кэш метрик памяти каждые 3 секунды и сохраняет в историю если активна запись */
func updateMemoryMetrics() {
	if usage, total, used, err := getmetrics.UsageMemory(); err == nil {
		cacheMutex.Lock()
		memCache = memoryPayload{
			MemoryUsage: usage,
			UsedMB:      used,
			TotalMemory: total,
			Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		}
		cacheMutex.Unlock()

		activeRecording, session := services.GetRecordingStatus()
		if activeRecording && session != nil {
			cacheMutex.RLock()
			cpuVal := cpuCache.CPU
			cacheMutex.RUnlock()

			shouldSave := cpuVal > session.CPUThreshold && usage > session.RAMThreshold
			if shouldSave {
				if err := saveMetricsHistory(cpuVal, usage, used, total); err != nil {
					log.Printf("Ошибка сохранения истории метрик: %v", err)
				}
			}
		}

		monitoringMutex.RLock()
		enabled := monitoringEnabled
		monitoringMutex.RUnlock()
		if enabled {
			cacheMutex.RLock()
			cpuVal := cpuCache.CPU
			cacheMutex.RUnlock()
			services.CheckAlerts(cpuVal, usage)
		}
	} else {
		log.Printf("Ошибка обновления кэша памяти: %v", err)
	}
}

/* Сохраняет метрики в историю в отдельной горутине чтобы не блокировать обновление кэша */
func saveMetricsHistory(cpuPercent, memoryPercent float64, memoryUsedMB, memoryTotalMB uint64) error {
	go func() {
		if err := services.SaveMetricsHistory(cpuPercent, memoryPercent, memoryUsedMB, memoryTotalMB); err != nil {
			log.Printf("Ошибка сохранения истории метрик: %v", err)
		}
	}()
	return nil
}

/* Обновляет кэш списка процессов каждые 5 секунд */
func updateProcessMetrics() {

	allConnections, err := net.Connections("all")
	if err != nil {
		log.Printf("Ошибка получения сетевых соединений: %v", err)
		allConnections = []net.ConnectionStat{}
	}

	if procs, err := getmetrics.UsageProcess(allConnections); err == nil {
		cacheMutex.Lock()
		procsCache = procs
		cacheMutex.Unlock()
	} else {
		log.Printf("Ошибка обновления кэша процессов: %v", err)
	}
}

/* Устанавливает ws соединение и начинает потоковую передачу метрик cpu каждую секунду */
func StreamCPU(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления соединения до WebSocket (CPU): %v", err)
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if enabled {

		if err := writeCPU(conn); err != nil {
			log.Printf("Ошибка отправки первого сообщения CPU: %v", err)
			return
		}
	} else {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(statusMsg)); err != nil {
			log.Printf("Ошибка отправки статуса мониторинга (CPU): %v", err)
			return
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := writeCPU(conn); err != nil {
			return
		}
	}
}

/* Отправляет кэшированные метрики cpu через ws соединение */
/* Проверяет состояние мониторинга перед отправкой данных */
func writeCPU(conn *websocket.Conn) error {

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if !enabled {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, []byte(statusMsg))
	}

	cacheMutex.RLock()
	data := cpuCache
	cacheMutex.RUnlock()

	if data.Timestamp == "" {
		statusMsg := `{"monitoringEnabled":true,"message":"Данные собираются...","cpu":0,"timestamp":"` + time.Now().Format("2006-01-02 15:04:05") + `"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, []byte(statusMsg))
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка сериализации метрик CPU: %v", err)
		return conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Ошибка сериализации данных"}`))
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, b)
}

/* Устанавливает ws соединение и начинает потоковую передачу метрик памяти каждые 3 секунды */
func StreamMemory(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления соединения до WebSocket (память): %v", err)
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if enabled {

		if err := writeMemory(conn); err != nil {
			log.Printf("Ошибка отправки первого сообщения памяти: %v", err)
			return
		}
	} else {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(statusMsg)); err != nil {
			log.Printf("Ошибка отправки статуса мониторинга (память): %v", err)
			return
		}
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := writeMemory(conn); err != nil {
			return
		}
	}
}

/* Отправляет кэшированные метрики памяти через ws соединение */
/* Проверяет состояние мониторинга перед отправкой данных */
func writeMemory(conn *websocket.Conn) error {

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if !enabled {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, []byte(statusMsg))
	}

	cacheMutex.RLock()
	data := memCache
	cacheMutex.RUnlock()

	if data.Timestamp == "" {
		statusMsg := `{"monitoringEnabled":true,"message":"Данные собираются...","memoryUsage":0,"usedMB":0,"totalMemory":0,"timestamp":"` + time.Now().Format("2006-01-02 15:04:05") + `"}`
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, []byte(statusMsg))
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка сериализации метрик памяти: %v", err)
		return conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Ошибка сериализации данных"}`))
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, b)
}

/* Устанавливает ws соединение и начинает потоковую передачу списка процессов каждые 5 секунд */
func StreamProcesses(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления соединения до WebSocket (процессы): %v", err)
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if enabled {

		if err := writeProcesses(conn); err != nil {
			log.Printf("Ошибка отправки первого сообщения процессов: %v", err)
			return
		}
	} else {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(statusMsg)); err != nil {
			log.Printf("Ошибка отправки статуса мониторинга (процессы): %v", err)
			return
		}
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := writeProcesses(conn); err != nil {
			return
		}
	}
}

/* Отправляет кэшированный список процессов через ws соединение */
/* Проверяет состояние мониторинга перед отправкой данных */
func writeProcesses(conn *websocket.Conn) error {

	monitoringMutex.RLock()
	enabled := monitoringEnabled
	monitoringMutex.RUnlock()

	if !enabled {

		statusMsg := `{"monitoringEnabled":false,"message":"Мониторинг выключен"}`
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, []byte(statusMsg))
	}

	cacheMutex.RLock()
	data := procsCache
	cacheMutex.RUnlock()

	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка сериализации списка процессов: %v", err)
		return conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"Ошибка сериализации данных"}`))
	}

	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, b)
}

/* Устанавливает состояние мониторинга и запускает или останавливает цикл обновления кэша */
func SetMonitoringEnabled(enabled bool) {
	monitoringMutex.Lock()
	wasEnabled := monitoringEnabled
	monitoringEnabled = enabled
	monitoringMutex.Unlock()

	if enabled && !wasEnabled {

		log.Println("Мониторинг включен: начинается сбор метрик")

		cacheCtx, cacheCancel = context.WithCancel(context.Background())
		go updateCacheLoop(cacheCtx)
	} else if !enabled && wasEnabled {

		log.Println("Мониторинг выключен: сбор метрик остановлен")
		if cacheCancel != nil {
			cacheCancel()
		}
	}
}

/* Возвращает текущее состояние мониторинга */
func GetMonitoringEnabled() bool {
	monitoringMutex.RLock()
	defer monitoringMutex.RUnlock()
	return monitoringEnabled
}
