package ui

import (
	"fmt"

	"github.com/RZhurakovskiy/agent/utils"
)

func ShowMainMenu() int {
	var action int

	fmt.Println("Главное меню:")
	menu := []MenuItem{
		{1, "Мониторинг и управление по загрузке CPU"},
		{2, "Мониторинг и управление по использованию памяти"},
		{3, "Запуск сервера и GUI"},
		{4, "Включить/выключить мониторинг системы"},
		{0, "Выйти"},
	}

	for _, item := range menu {
		fmt.Printf(" [%d] %s\n", item.ID, item.Text)
	}
	fmt.Println("---------------------------------")
	action = utils.GetUserInput()
	return action
}
