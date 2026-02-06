package utils

import "fmt"

func GetUserInput() int {
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Ошибка! Введите номер пункта (цифру).")

		ClearScanBuffer()
		return -1
	}
	return choice
}
