package utils

import "fmt"

func ClearScanBuffer() {
	var clear string
	fmt.Scanln(&clear)
}
