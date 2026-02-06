package cpu

import (
	"fmt"

	"github.com/RZhurakovskiy/agent/ui"
	"github.com/RZhurakovskiy/agent/utils"
)

func ProcessMenu() {
	for {
		choice := ui.ViewMenu()

		switch choice {
		case 1:
			viewProcess()

		case 2:
			action, pid, name := ui.CompletionMenu()
			switch action {
			case 1:
				killProcessByPIDOrName(pid, "")
			case 2:
				killProcessByPIDOrName(0, name)
			case 3:
				filteredProcess(pid, "")
			case 4:
				filteredProcess(0, name)
			case 0:

			default:
				fmt.Println("Неверный выбор в подменю.")
			}

		case 3:
			checkingSuspiciousActivity()
			action := ui.СheckSuspiciousActivityMenu()

			switch action {
			case 0:

			}
		case 4:
			var thresHold float64
			fmt.Print("Процессы с загрузкой CPU выше этого значения будут записаны в лог. Укажите порог (%): ")
			_, err := fmt.Scan(&thresHold)
			if err != nil {
				fmt.Println("Ошибка! Используйте число (например, 70.5 или 70).")
				utils.ClearScanBuffer()
			} else {
				startDaemonMode(thresHold)
			}

		case 5:
			var pid int32
			fmt.Print("Введите PID процесса для получения статистики: ")
			_, err := fmt.Scan(&pid)
			if err != nil {
				fmt.Println("Ошибка! Введите корректный PID.")
				utils.ClearScanBuffer()
			} else {
				PrintProcessStats(pid)
			}

		case 6:
			exportProcessesMenu()

		case 7:
			StartRecordingMenu()

		case 0:
			fmt.Println("\nВыход из программы.")
			return

		default:
			fmt.Println("Неверный пункт меню!")
		}
	}
}
