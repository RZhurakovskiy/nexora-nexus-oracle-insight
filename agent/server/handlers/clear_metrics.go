package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RZhurakovskiy/agent/server/services"
)

func ClearMetrics(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	if err := services.ClearMetricsHistory(); err != nil {
		http.Error(writer, "Ошибка очистки метрик: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"success": true,
		"message": "Метрики успешно очищены",
	})
}


