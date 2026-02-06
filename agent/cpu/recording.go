package cpu

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RZhurakovskiy/agent/server/services"
	"github.com/RZhurakovskiy/agent/server/ws"
	_ "github.com/mattn/go-sqlite3"
)

func StartRecordingMenu() {
	fmt.Println("\n=== Запись метрик процессов ===")
	fmt.Println("Эта функция записывает процессы, превышающие пороги CPU и RAM")
	fmt.Println("в базу данных для последующего анализа.")
	fmt.Println()

	var cpuThreshold float64
	var ramThreshold float64
	var duration int

	fmt.Print("Введите порог загрузки CPU (%): ")
	_, err := fmt.Scan(&cpuThreshold)
	if err != nil {
		fmt.Println("Ошибка! Используйте число (например, 70.5 или 70).")
		return
	}

	if cpuThreshold <= 0 || cpuThreshold > 100 {
		fmt.Println("Порог CPU должен быть от 0 до 100%")
		return
	}

	fmt.Print("Введите порог использования RAM (%): ")
	_, err = fmt.Scan(&ramThreshold)
	if err != nil {
		fmt.Println("Ошибка! Используйте число (например, 80.5 или 80).")
		return
	}

	if ramThreshold <= 0 || ramThreshold > 100 {
		fmt.Println("Порог RAM должен быть от 0 до 100%")
		return
	}

	fmt.Print("Введите продолжительность записи в секундах (минимум 60): ")
	_, err = fmt.Scan(&duration)
	if err != nil {
		fmt.Println("Ошибка! Используйте число.")
		return
	}

	if duration < 60 {
		fmt.Println("Продолжительность должна быть не менее 60 секунд")
		return
	}

	fmt.Println("\nИнициализация базы данных...")
	sqlDB, err := sql.Open("sqlite3", "./monitor.db")
	if err != nil {
		fmt.Printf("Ошибка открытия базы данных: %v\n", err)
		return
	}
	defer sqlDB.Close()

	services.SetDB(sqlDB)

	fmt.Println("Включение мониторинга...")
	ws.SetMonitoringEnabled(true)

	fmt.Println("Запуск записи метрик...")
	sessionID, err := services.StartRecording(cpuThreshold, ramThreshold, duration)
	if err != nil {
		fmt.Printf("Ошибка запуска записи: %v\n", err)
		ws.SetMonitoringEnabled(false)
		return
	}

	fmt.Printf("\n✓ Запись метрик запущена!\n")
	fmt.Printf("  ID сессии: %d\n", sessionID)
	fmt.Printf("  Порог CPU: %.2f%%\n", cpuThreshold)
	fmt.Printf("  Порог RAM: %.2f%%\n", ramThreshold)
	fmt.Printf("  Продолжительность: %d секунд\n", duration)
	fmt.Printf("  Завершится в: %s\n", time.Now().Add(time.Duration(duration)*time.Second).Format("2006-01-02 15:04:05"))
	fmt.Println("\nЗапись будет автоматически остановлена по истечении времени.")
	fmt.Println("Для досрочной остановки нажмите Ctrl+C.")
	fmt.Println("----------------------------------------")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	timer := time.NewTimer(time.Duration(duration) * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("\nВремя записи истекло. Остановка записи...")
		services.StopRecording()
		ws.SetMonitoringEnabled(false)
		fmt.Println("✓ Запись завершена. Мониторинг выключен.")
	case <-sigChan:
		fmt.Println("\n\nПолучен сигнал прерывания. Остановка записи...")
		services.StopRecording()
		ws.SetMonitoringEnabled(false)
		fmt.Println("✓ Запись остановлена. Мониторинг выключен.")
	}
}
