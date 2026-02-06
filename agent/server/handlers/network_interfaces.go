package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
)

func GetNetworkInterfaces(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	result, err := getmetrics.GetNetworkInterfacesIO()
	if err != nil {
		http.Error(writer, "Ошибка получения сетевых интерфейсов", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(result)
}


