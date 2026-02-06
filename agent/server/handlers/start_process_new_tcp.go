package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/RZhurakovskiy/agent/server/services"
)

func StartProcess(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Только POST разрешён", http.StatusMethodNotAllowed)
		return
	}

	var req models.StartProcessRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "Неверный JSON", http.StatusBadRequest)
		return
	}

	result, err := services.StartProcess(req.Command, req.Args, req.Cwd)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(models.StartProcessResponse{
		PID:     result.PID,
		Command: req.Command,
		Args:    req.Args,
		Cwd:     req.Cwd,
		Msg:     result.Msg,
	})
}
