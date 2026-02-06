package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/process"
)

func processExists(pid int32) (bool, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return false, nil
	}

	_, err = proc.Status()
	if err != nil {
		return false, nil
	}

	return true, nil
}

func KillProcessById(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Метод не разрешён. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		PID int32 `json:"pid"`
	}

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		log.Printf("Ошибка декодирования JSON в KillProcessById: %v", err)
		http.Error(writer, "Некорректный JSON. Ожидается: {\"pid\": <число>}", http.StatusBadRequest)
		return
	}

	pid := input.PID

	if pid <= 0 {
		http.Error(writer, "Некорректный PID. PID должен быть положительным числом", http.StatusBadRequest)
		return
	}

	var message string

	exists, err := processExists(pid)
	if err != nil {
		log.Printf("Ошибка проверки существования процесса %d: %v", pid, err)
		message = "Ошибка проверки процесса: " + err.Error()
	} else if !exists {
		message = "Процесс не найден"
	} else {

		proc, err := process.NewProcess(pid)
		if err != nil {
			log.Printf("Ошибка создания объекта процесса %d: %v", pid, err)
			message = "Не удалось загрузить процесс: " + err.Error()
		} else {

			err = proc.Kill()
			if err != nil {
				log.Printf("Ошибка завершения процесса %d: %v", pid, err)
				message = "Не удалось завершить процесс: " + err.Error()
			} else {
				log.Printf("Процесс %d успешно завершен", pid)
				message = "Процесс успешно завершен"
			}
		}
	}

	response := models.KillProcessByID{
		PID:       pid,
		Message:   message,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Ошибка сериализации ответа в KillProcessById: %v", err)
		http.Error(writer, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}
}
