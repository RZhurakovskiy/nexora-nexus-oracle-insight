package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
)

func GetNetworkTopProcesses(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	limit := 20
	if v := request.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	result, err := getmetrics.GetTopNetworkProcesses(limit)
	if err != nil {
		http.Error(writer, "Ошибка получения сетевой статистики", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(result)
}


