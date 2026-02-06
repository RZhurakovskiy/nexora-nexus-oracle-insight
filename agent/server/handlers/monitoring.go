// Обработчики для управления состоянием мониторинга системы
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/RZhurakovskiy/agent/server/ws"
)

/* Устанавливает состояние мониторинга через API и возвращает результат операции */
func SetMonitoringStatus(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	var input models.MonitoringStatusRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		log.Printf("Ошибка декодирования JSON в SetMonitoringStatus: %v", err)
		http.Error(writer, "Некорректный JSON. Ожидается: {\"enabled\": <true|false>}", http.StatusBadRequest)
		return
	}

	wasEnabled := ws.GetMonitoringEnabled()

	ws.SetMonitoringEnabled(input.Enabled)

	currentEnabled := ws.GetMonitoringEnabled()

	var message string
	if currentEnabled {
		if wasEnabled {
			message = "Мониторинг уже был включен"
		} else {
			message = "Мониторинг успешно включен"
			log.Println("Мониторинг включен через API")
		}
	} else {
		if !wasEnabled {
			message = "Мониторинг уже был выключен"
		} else {
			message = "Мониторинг успешно выключен"
			log.Println("Мониторинг выключен через API")
		}
	}

	response := models.MonitoringStatusResponse{
		Enabled:   currentEnabled,
		Message:   message,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Ошибка сериализации ответа в SetMonitoringStatus: %v", err)
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}

/* Возвращает текущее состояние мониторинга в формате json */
func GetMonitoringStatus(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	enabled := ws.GetMonitoringEnabled()

	message := "Мониторинг выключен"
	if enabled {
		message = "Мониторинг включен"
	}

	response := models.MonitoringStatusResponse{
		Enabled:   enabled,
		Message:   message,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Ошибка сериализации ответа в GetMonitoringStatus: %v", err)
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}
