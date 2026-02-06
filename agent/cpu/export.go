package cpu

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RZhurakovskiy/agent/server/getmetrics"
	"github.com/shirou/gopsutil/v4/net"
)

func exportProcessesMenu() {
	fmt.Println("\nЭкспорт процессов:")
	fmt.Println(" [1] Экспорт в JSON")
	fmt.Println(" [2] Экспорт в CSV")
	fmt.Println(" [0] Отмена")
	fmt.Print("Выберите формат: ")

	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Ошибка ввода!")
		return
	}

	switch choice {
	case 1:
		exportProcessesJSON()
	case 2:
		exportProcessesCSV()
	case 0:
		return
	default:
		fmt.Println("Неверный выбор!")
	}
}

func exportProcessesJSON() {
	fmt.Println("\nПолучение списка процессов...")

	allConnections, err := net.Connections("all")
	if err != nil {
		log.Printf("Ошибка получения сетевых соединений: %v", err)
		allConnections = []net.ConnectionStat{}
	}

	processes, err := getmetrics.UsageProcess(allConnections)
	if err != nil {
		fmt.Printf("Ошибка получения процессов: %v\n", err)
		return
	}

	filename := fmt.Sprintf("processes_%s.json", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(processes); err != nil {
		fmt.Printf("Ошибка записи JSON: %v\n", err)
		return
	}

	fmt.Printf("Процессы успешно экспортированы в файл: %s\n", filename)
	fmt.Printf("Всего процессов: %d\n", len(processes))
}

func exportProcessesCSV() {
	fmt.Println("\nПолучение списка процессов...")

	allConnections, err := net.Connections("all")
	if err != nil {
		log.Printf("Ошибка получения сетевых соединений: %v", err)
		allConnections = []net.ConnectionStat{}
	}

	processes, err := getmetrics.UsageProcess(allConnections)
	if err != nil {
		fmt.Printf("Ошибка получения процессов: %v\n", err)
		return
	}

	filename := fmt.Sprintf("processes_%s.csv", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"PID", "Имя", "Путь", "Командная строка", "Пользователь",
		"Статус", "Время создания", "Родительский PID",
		"CPU %", "Память %", "Память RSS (байты)", "Порты",
	}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Ошибка записи заголовков CSV: %v\n", err)
		return
	}

	for _, p := range processes {
		portsStr := ""
		if len(p.Ports) > 0 {
			portStrs := make([]string, len(p.Ports))
			for i, port := range p.Ports {
				portStrs[i] = strconv.FormatUint(uint64(port), 10)
			}
			portsStr = fmt.Sprintf("[%s]", fmt.Sprint(portStrs))
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
		if err := writer.Write(record); err != nil {
			fmt.Printf("Ошибка записи строки CSV: %v\n", err)
			return
		}
	}

	fmt.Printf("Процессы успешно экспортированы в файл: %s\n", filename)
	fmt.Printf("Всего процессов: %d\n", len(processes))
}
