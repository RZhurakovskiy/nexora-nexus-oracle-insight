package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
)

func GetDiskHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	resp := getmetrics.GetDiskHealth(runtime.GOOS)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}



