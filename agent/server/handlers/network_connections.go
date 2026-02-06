package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
)

func GetNetworkConnections(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешен. Разрешён только GET", http.StatusMethodNotAllowed)
		return
	}

	kind := request.URL.Query().Get("kind")
	result, err := getmetrics.GetNetworkConnections(kind)
	if err != nil {
		http.Error(writer, "Ошибка получения сетевых соединений", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(result)
}


