package cpu

import (
	"database/sql"
	"fmt"

	"github.com/RZhurakovskiy/agent/server/services"
	"github.com/RZhurakovskiy/agent/server/ws"
	_ "github.com/mattn/go-sqlite3"
)

func ToggleMonitoringMenu() {
	fmt.Println("\n=== Управление мониторингом системы ===")

	sqlDB, err := sql.Open("sqlite3", "./monitor.db")
	if err != nil {
		fmt.Printf("Ошибка открытия базы данных: %v\n", err)
		return
	}
	defer sqlDB.Close()

	services.SetDB(sqlDB)

	currentStatus := ws.GetMonitoringEnabled()

	if currentStatus {
		fmt.Println("Текущий статус: Мониторинг ВКЛЮЧЕН")
		fmt.Print("Выключить мониторинг? (y/n): ")
		var answer string
		fmt.Scan(&answer)
		if answer == "y" || answer == "Y" {
			ws.SetMonitoringEnabled(false)
			fmt.Println("✓ Мониторинг выключен")
		} else {
			fmt.Println("Отменено")
		}
	} else {
		fmt.Println("Текущий статус: Мониторинг ВЫКЛЮЧЕН")
		fmt.Print("Включить мониторинг? (y/n): ")
		var answer string
		fmt.Scan(&answer)
		if answer == "y" || answer == "Y" {
			ws.SetMonitoringEnabled(true)
			fmt.Println("✓ Мониторинг включен")
		} else {
			fmt.Println("Отменено")
		}
	}
}
