/* Обработчики для работы с алертами системы */
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RZhurakovskiy/agent/server/services"
)

/* Получает список алертов с поддержкой фильтрации по неподтвержденным и лимиту записей
   Возвращает JSON массив с алертами */
func GetAlerts(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	limitStr := request.URL.Query().Get("limit")
	unacknowledgedOnlyStr := request.URL.Query().Get("unacknowledged_only")

	limit := 100
	unacknowledgedOnly := false

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if unacknowledgedOnlyStr == "true" || unacknowledgedOnlyStr == "1" {
		unacknowledgedOnly = true
	}

	alerts, err := services.GetAlerts(limit, unacknowledgedOnly)
	if err != nil {
		http.Error(writer, "Ошибка получения алертов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(alerts); err != nil {
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}

/* Отмечает алерт как подтвержденный по его ID */
func AcknowledgeAlert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	if err := services.AcknowledgeAlert(req.ID); err != nil {
		http.Error(writer, "Ошибка подтверждения алерта: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"success": true,
		"message": "Алерт подтвержден",
	})
}

/* Устанавливает пороги для CPU и памяти, при превышении которых будут создаваться алерты */
func SetAlertThresholds(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CPUThreshold    float64 `json:"cpuThreshold"`
		MemoryThreshold float64 `json:"memoryThreshold"`
	}

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	services.SetAlertThresholds(req.CPUThreshold, req.MemoryThreshold)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"success":        true,
		"cpuThreshold":   req.CPUThreshold,
		"memoryThreshold": req.MemoryThreshold,
		"message":        "Пороги установлены",
	})
}

/* Возвращает текущие пороги для CPU и памяти в формате JSON */
func GetAlertThresholds(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	thresholds := services.GetAlertThresholds()

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(thresholds)
}


