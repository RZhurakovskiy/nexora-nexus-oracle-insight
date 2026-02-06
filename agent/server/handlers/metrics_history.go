/* Обработчики для работы с историей метрик системы */
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/RZhurakovskiy/agent/server/services"
)

/* Получает историю метрик за указанный период времени с поддержкой фильтрации по датам и лимиту записей
   Возвращает JSON массив с историей метрик */
func GetMetricsHistory(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	fromStr := request.URL.Query().Get("from")
	toStr := request.URL.Query().Get("to")
	limitStr := request.URL.Query().Get("limit")

	now := time.Now()
	from := now.AddDate(0, 0, -7)
	to := now
	limit := 1000

	if fromStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", fromStr); err == nil {
			from = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", fromStr); err == nil {
			from = parsed
		} else if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = parsed
		} else {
			http.Error(writer, "Некорректный формат даты 'from': "+fromStr, http.StatusBadRequest)
			return
		}
	}

	if toStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", toStr); err == nil {
			to = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", toStr); err == nil {
			to = parsed
		} else if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			to = parsed
		} else {
			http.Error(writer, "Некорректный формат даты 'to': "+toStr, http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	history, err := services.GetMetricsHistory(from, to, limit)
	if err != nil {
		http.Error(writer, "Ошибка получения истории метрик: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(history); err != nil {
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}
