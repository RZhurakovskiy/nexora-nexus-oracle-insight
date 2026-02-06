package ui

import (
	"fmt"

	"github.com/RZhurakovskiy/agent/utils"
)

func ViewMenu() int {
	menu := []MenuItem{
		{1, "Посмотреть список процессов"},
		{2, "Завершить или найти процесс по PID или названию"},
		{3, "Проверить подозрительные процессы"},
		{4, "Запустить фоновый мониторинг загрузки CPU"},
		{5, "Статистика процесса по PID"},
		{6, "Экспорт процессов в CSV/JSON"},
		{7, "Запись метрик процессов (с порогами CPU/RAM)"},
		{0, "Выйти"},
	}
	fmt.Println("Главное меню администрирования процессами:")
	for _, item := range menu {
		fmt.Printf(" [%d] %s\n", item.ID, item.Text)
	}
	fmt.Println("---------------------------------")

	return getUserInput()
}

func CompletionMenu() (int, int32, string) {
	var action int
	var pid int32
	var name string
	menu := []MenuItem{
		{1, "Завершить процесс по PID"},
		{2, "Завершить процесс по названию"},
		{3, "Найти процессы по PID"},
		{4, "Найти процессы по названию"},
		{0, "Вернуться в главное меню"},
	}

	fmt.Println("\nДополнительное меню:")
	for _, item := range menu {
		fmt.Printf(" [%d] %s\n", item.ID, item.Text)
	}

	action = getUserInput()

	switch action {
	case 1:
		fmt.Print("\nВведите PID процесса для завершения: ")
		pid = readPID()
	case 2:
		fmt.Print("\nВведите название процесса для завершения: ")
		name = readName()
	case 3:
		fmt.Print("\nВведите PID процесса для поиска: ")
		pid = readPID()
	case 4:
		fmt.Print("\nВведите название процесса для поиска: ")
		name = readName()
	}

	return action, pid, name
}

func СheckSuspiciousActivityMenu() int {
	var action int
	menu := []MenuItem{
		{0, "Вернуться в главное меню"},
	}

	for _, item := range menu {
		fmt.Printf(" [%d] %s\n", item.ID, item.Text)
	}

	action = getUserInput()
	return action
}

func getUserInput() int {
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Ошибка! Введите номер пункта (цифру).")
		utils.ClearScanBuffer()
		return -1
	}
	return choice
}

func readPID() int32 {
	var pid int32
	_, err := fmt.Scan(&pid)
	if err != nil {
		fmt.Println("Неверный PID!")
		utils.ClearScanBuffer()
		return 0
	}
	if pid <= 0 {
		fmt.Println("PID должен быть положительным числом.")
		return 0
	}
	return pid
}

func readName() string {
	var name string
	_, err := fmt.Scan(&name)
	if err != nil {
		fmt.Println("Ошибка ввода названия!")
		utils.ClearScanBuffer()
		return ""
	}
	return name
}
