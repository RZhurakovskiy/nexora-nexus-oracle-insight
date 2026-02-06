/* Точка входа в приложение агента мониторинга системы */
package main

import (
	"fmt"

	"github.com/RZhurakovskiy/agent/cpu"
	"github.com/RZhurakovskiy/agent/server/api"

	"github.com/RZhurakovskiy/agent/ui"
)

/* Запускает главное меню и обрабатывает выбор пользователя */
func main() {
	ui.ShowBanner()

	action := ui.ShowMainMenu()

	switch action {
	case 1:
		cpu.ProcessMenu()
	case 2:
		fmt.Println("Раздел в разработке")
		return
	case 3:
		port := "8080"
		api.StartServer(port)

		return
	case 4:
		cpu.ToggleMonitoringMenu()
	case 0:
		fmt.Println("Выход...")
		return
	}

}
