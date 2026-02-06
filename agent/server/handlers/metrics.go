package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/user"
)

func GetHostUserName(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	userName, userErr := user.Current()
	hostName, hostErr := os.Hostname()

	if userErr != nil {
		http.Error(writer, "Не удалось получить имя пользователя", http.StatusInternalServerError)
		return
	}
	if hostErr != nil {
		http.Error(writer, "Не удалось получить имя хоста", http.StatusInternalServerError)
		return
	}

	type SystemInfo struct {
		Username string `json:"username"`
		Hostname string `json:"hostname"`
	}

	info := SystemInfo{
		Username: userName.Username,
		Hostname: hostName,
	}

	writer.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(writer)
	err := encoder.Encode(info)
	if err != nil {
		log.Printf("Ошибка сериализации информации о системе: %v", err)
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}
