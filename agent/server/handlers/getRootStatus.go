package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/RZhurakovskiy/agent/server/models"
)

func GetRootStatus(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	var rootStatus bool
	if os.Geteuid() == 0 {
		rootStatus = true
	} else {
		rootStatus = false
	}

	resp := models.RoorStatus{
		RootStatus: rootStatus,
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)

}
