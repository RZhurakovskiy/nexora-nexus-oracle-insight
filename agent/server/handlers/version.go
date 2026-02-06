/* Обработчики для получения версии сервера */
package handlers

import (
	"encoding/json"
	"net/http"
)

const ServerVersion = "1.1.0"

/* Структура для ответа с версией сервера */
type VersionResponse struct {
	ServerVersion string `json:"serverVersion"`
}

/* Возвращает версию сервера в формате JSON */
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(writer).Encode(VersionResponse{
		ServerVersion: ServerVersion,
	})
}
