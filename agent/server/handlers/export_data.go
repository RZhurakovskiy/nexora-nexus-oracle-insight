/* Обработчики для экспорта данных в CSV и JSON форматы */
package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
)

/* Экспортирует список всех процессов в csv или json формат в зависимости от параметра format */
func ExportProcesses(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	format := request.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	allConnections, err := net.Connections("all")
	if err != nil {
		http.Error(writer, "Ошибка получения сетевых соединений", http.StatusInternalServerError)
		return
	}

	processes, err := getmetrics.UsageProcess(allConnections)
	if err != nil {
		http.Error(writer, "Ошибка получения процессов", http.StatusInternalServerError)
		return
	}

	switch strings.ToLower(format) {
	case "csv":
		exportProcessesCSV(writer, processes)
	case "json":
		exportProcessesJSON(writer, processes)
	default:
		http.Error(writer, "Неподдерживаемый формат. Используйте 'csv' или 'json'", http.StatusBadRequest)
	}
}

/* Записывает список процессов в CSV формат с заголовками и данными */
func exportProcessesCSV(writer http.ResponseWriter, processes []models.ProcessInfo) {
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=processes_%s.csv", time.Now().Format("20060102_150405")))

	w := csv.NewWriter(writer)
	defer w.Flush()

	headers := []string{
		"PID", "Имя", "Путь", "Командная строка", "Пользователь",
		"Статус", "Время создания", "Родительский PID",
		"CPU %", "Память %", "Память RSS (байты)", "Порты",
	}
	if err := w.Write(headers); err != nil {
		log.Printf("Ошибка записи CSV заголовков: %v", err)
		return
	}

	for _, p := range processes {
		portsStr := ""
		if len(p.Ports) > 0 {
			portStrs := make([]string, len(p.Ports))
			for i, port := range p.Ports {
				portStrs[i] = strconv.FormatUint(uint64(port), 10)
			}
			portsStr = strings.Join(portStrs, ";")
		}

		createTimeStr := ""
		if p.CreateTime > 0 {
			createTimeStr = time.Unix(p.CreateTime/1000, 0).Format("2006-01-02 15:04:05")
		}

		record := []string{
			strconv.FormatInt(int64(p.PID), 10),
			p.Name,
			p.Exe,
			p.Cmdline,
			p.Username,
			p.Status,
			createTimeStr,
			strconv.FormatInt(int64(p.ParentPID), 10),
			fmt.Sprintf("%.2f", p.CPUPercent),
			fmt.Sprintf("%.2f", p.MemoryPercent),
			strconv.FormatUint(p.MemoryRSS, 10),
			portsStr,
		}
		if err := w.Write(record); err != nil {
			log.Printf("Ошибка записи CSV строки: %v", err)
			return
		}
	}
}

/* Записывает список процессов в JSON формат с отступами для читаемости */
func exportProcessesJSON(writer http.ResponseWriter, processes []models.ProcessInfo) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=processes_%s.json", time.Now().Format("20060102_150405")))

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(processes); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
		return
	}
}

/* Экспортирует текущие метрики CPU и памяти в CSV или JSON формат */
func ExportMetrics(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	format := request.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	cpuUsage, err := getmetrics.UsageCPU(100 * time.Millisecond)
	if err != nil {
		http.Error(writer, "Ошибка получения метрик CPU", http.StatusInternalServerError)
		return
	}

	memUsage, totalMem, usedMem, err := getmetrics.UsageMemory()
	if err != nil {
		http.Error(writer, "Ошибка получения метрик памяти", http.StatusInternalServerError)
		return
	}

	metrics := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"cpu": map[string]interface{}{
			"percent": cpuUsage,
		},
		"memory": map[string]interface{}{
			"percent":    memUsage,
			"usedMB":     usedMem,
			"totalMB":    totalMem,
			"usedBytes":  usedMem * 1024 * 1024,
			"totalBytes": totalMem * 1024 * 1024,
		},
	}

	switch strings.ToLower(format) {
	case "csv":
		exportMetricsCSV(writer, metrics)
	case "json":
		exportMetricsJSON(writer, metrics)
	default:
		http.Error(writer, "Неподдерживаемый формат. Используйте 'csv' или 'json'", http.StatusBadRequest)
	}
}

/* Записывает метрики в CSV формат с одной строкой данных */
func exportMetricsCSV(writer http.ResponseWriter, metrics map[string]interface{}) {
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=metrics_%s.csv", time.Now().Format("20060102_150405")))

	w := csv.NewWriter(writer)
	defer w.Flush()

	headers := []string{
		"Время", "CPU %", "Память %", "Память использовано (МБ)", "Память всего (МБ)",
		"Память использовано (байты)", "Память всего (байты)",
	}
	if err := w.Write(headers); err != nil {
		log.Printf("Ошибка записи CSV заголовков: %v", err)
		return
	}

	cpuData := metrics["cpu"].(map[string]interface{})
	memData := metrics["memory"].(map[string]interface{})

	record := []string{
		metrics["timestamp"].(string),
		fmt.Sprintf("%.2f", cpuData["percent"]),
		fmt.Sprintf("%.2f", memData["percent"]),
		fmt.Sprintf("%d", memData["usedMB"]),
		fmt.Sprintf("%d", memData["totalMB"]),
		fmt.Sprintf("%d", memData["usedBytes"]),
		fmt.Sprintf("%d", memData["totalBytes"]),
	}
	if err := w.Write(record); err != nil {
		log.Printf("Ошибка записи CSV строки: %v", err)
		return
	}
}

/* Записывает метрики в JSON формат с отступами */
func exportMetricsJSON(writer http.ResponseWriter, metrics map[string]interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=metrics_%s.json", time.Now().Format("20060102_150405")))

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metrics); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
		return
	}
}
