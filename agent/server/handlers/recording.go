/* Обработчики для управления записью процессов с высоким использованием ресурсов */
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RZhurakovskiy/agent/server/services"
	"github.com/RZhurakovskiy/agent/server/ws"
)

/* Структура для запроса на начало записи процессов */
type StartRecordingRequest struct {
	CPUThreshold float64 `json:"cpuThreshold"`
	RAMThreshold float64 `json:"ramThreshold"`
	Duration     int     `json:"duration"`
}

/* Запускает новую сессию записи процессов с указанными порогами cpu и ram */
func StartRecording(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	if !ws.GetMonitoringEnabled() {
		http.Error(writer, "Для запуска записи метрик необходимо включить мониторинг", http.StatusBadRequest)
		return
	}

	var req StartRecordingRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	if req.CPUThreshold <= 0 || req.RAMThreshold <= 0 || req.Duration <= 0 {
		http.Error(writer, "Некорректные параметры", http.StatusBadRequest)
		return
	}

	sessionID, err := services.StartRecording(req.CPUThreshold, req.RAMThreshold, req.Duration)
	fmt.Println(sessionID)
	if err != nil {
		http.Error(writer, "Ошибка запуска записи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"success":   true,
		"sessionId": sessionID,
		"message":   "Запись метрик запущена",
	})
}

/* Останавливает активную сессию записи процессов */
func StopRecording(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	if err := services.StopRecording(); err != nil {
		http.Error(writer, "Ошибка остановки записи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"success": true,
		"message": "Запись метрик остановлена",
	})
}

/* Возвращает статус текущей сессии записи включая пороги и временные метки */
func GetRecordingStatus(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	active, session := services.GetRecordingStatus()

	response := map[string]interface{}{
		"active": active,
	}

	if session != nil {
		response["session"] = map[string]interface{}{
			"id":           session.ID,
			"cpuThreshold": session.CPUThreshold,
			"ramThreshold": session.RAMThreshold,
			"duration":     session.Duration,
			"startedAt":    session.StartedAt.Format(time.RFC3339),
			"endTime":      session.EndTime.Format(time.RFC3339),
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

/* Получает список записанных процессов для указанной сессии с поддержкой лимита */
func GetRecordedProcesses(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	sessionIDStr := request.URL.Query().Get("sessionId")
	limitStr := request.URL.Query().Get("limit")

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		http.Error(writer, "Некорректный sessionId", http.StatusBadRequest)
		return
	}

	limit := 1000
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	processes, err := services.GetRecordedProcesses(sessionID, limit)
	if err != nil {
		http.Error(writer, "Ошибка получения записанных процессов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(processes)
}
