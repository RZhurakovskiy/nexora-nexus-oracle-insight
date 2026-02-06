package cpu

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var CPUThresholdPercent float64

func viewProcess() {
	procs, err := process.Processes()
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Сканирование процессов... (измерение займёт ~1 секунду)")
	fmt.Println("-------------------------------------------------------")
	if err != nil {
		log.Println("Не удалось получить список процессов:", err)
		return
	}

	qntProcess := 0
	fmt.Printf("\n%-10s %-20s\n", "PID", "Название")
	fmt.Printf("%-10s %-20s\n", "----------", "--------------------")
	for _, proc := range procs {
		pid := proc.Pid
		name, err := proc.Name()
		if err != nil {
			name = "процесс неопределен"
		}

		qntProcess++
		fmt.Printf("%-10d %-20s\n", pid, name)
	}
	fmt.Println("---------------------------------")
	fmt.Printf("Найдено процессов: %d\n", qntProcess)
	fmt.Println("---------------------------------")
}

func killProcessByPIDOrName(processPID int32, processName string) {
	procs, err := process.Processes()
	if err != nil {
		log.Println("Не удалось получить список процессов:", err)
		return
	}

	foundByPID := false
	foundByName := false

	for _, proc := range procs {
		pid := proc.Pid
		name, err := proc.Name()
		if err != nil {
			name = ""
		}

		if processPID != 0 && processPID == pid {
			err := proc.Kill()
			if err != nil {
				fmt.Printf("Не удалось завершить процесс с PID %d: %v\n", pid, err)
			} else {
				fmt.Printf("Процесс с PID %d успешно завершён.\n", pid)
			}
			foundByPID = true
		}

		if processName != "" && processName == name {
			err := proc.Kill()
			if err != nil {
				fmt.Printf("Не удалось завершить процесс '%s' (PID %d): %v\n", name, pid, err)
			} else {
				fmt.Printf("Процесс '%s' (PID %d) успешно завершён.\n", name, pid)
			}
			foundByName = true
		}
	}

	if processPID != 0 && !foundByPID {
		fmt.Println("Процесс по PID не найден.")
	}

	if processName != "" && !foundByName {
		fmt.Println("Процесс по названию не найден.")
	}
}

func filteredProcess(processPIDSearch int32, processNameSearch string) {
	procs, err := process.Processes()
	if err != nil {
		log.Println("Не удалось получить список процессов:", err)
		return
	}

	qntProcess := 0
	fmt.Printf("%-10s %-20s\n", "PID", "Название")
	fmt.Printf("%-10s %-20s\n", "----------", "--------------------")

	for _, proc := range procs {
		pid := proc.Pid
		name, err := proc.Name()
		if err != nil {
			name = "процесс неопределен"
		}

		matched := false
		if processPIDSearch != 0 && pid == processPIDSearch {
			matched = true
		}
		if processNameSearch != "" && processNameSearch == name {
			matched = true
		}

		if matched {
			qntProcess++
			fmt.Printf("%-10d %-20s\n", pid, name)
		}
	}

	if qntProcess == 0 {
		fmt.Println("Совпадений не найдено.")
	} else {
		fmt.Printf("\nНайдено процессов: %d\n", qntProcess)
	}
}

func startDaemonMode(threshold float64) {
	CPUThresholdPercent = threshold

	fmt.Println("\nЗапуск фонового мониторинга...")
	fmt.Println("Проверка каждые 5 секунд. CPU >", threshold, "% → запись в лог.")
	fmt.Println("Для остановки нажмите Ctrl+C.")
	fmt.Println("----------------------------------------")

	logFile, err := os.OpenFile("логирование_загрузки.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Ошибка создания лог-файла: %v\n", err)
		return
	}
	defer logFile.Close()

	logFile.WriteString("=== Начало мониторинга ===\n")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkAndLogHeavyProcesses(logFile)
		case <-sigChan:
			fmt.Println("\nМониторинг остановлен. Возврат в меню...")
			logFile.WriteString("=== Мониторинг остановлен ===\n")
			return
		}
	}
}

func checkingSuspiciousActivity() {
	fmt.Println("Проверка подозрительных процессов: в разработке")
}

func checkAndLogHeavyProcesses(logFile *os.File) {
	procs, err := process.Processes()
	if err != nil {
		fmt.Printf("Ошибка получения процессов: %v\n", err)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "\n[%s] Проверка процессов:\n", now)

	heavyFound := false

	for _, proc := range procs {
		pid := proc.Pid
		name, _ := proc.Name()

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		if cpuPercent > CPUThresholdPercent {
			heavyFound = true
			fmt.Fprintf(logFile, "  PID: %d | CPU: %.1f%% | Название: %s\n", pid, cpuPercent, name)
			fmt.Printf("Обнаружен процесс с высокой загрузкой > %.1f%%: PID=%d, CPU=%.1f%%, %s\n",
				CPUThresholdPercent, pid, cpuPercent, name)
		}
	}

	if !heavyFound {
		fmt.Fprintf(logFile, "  Нет процессов с CPU > %.1f%%.\n", CPUThresholdPercent)
	}
}
