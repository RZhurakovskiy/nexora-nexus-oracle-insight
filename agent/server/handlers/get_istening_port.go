package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
)

func GetListeningPort(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешен. Разрешён только GET", http.StatusMethodNotAllowed)
		return
	}

	result, err := getmetrics.GetListeningPorts()
	if err != nil {

		http.Error(writer, "Ошибка получения информации о портах:", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(result)
}
