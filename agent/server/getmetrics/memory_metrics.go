/* Функции для получения метрик использования памяти */
package getmetrics

import (
	"github.com/shirou/gopsutil/v4/mem"
)

/* Получает метрики использования памяти включая процент использования, общий объем и использованный объем в мегабайтах */
/* Возвращает процент использования, общий объем в мб, использованный объем в мб или ошибку */
func UsageMemory() (float64, uint64, uint64, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0, err
	}

	const mib = uint64(1024 * 1024)
	totalMiB := memory.Total / mib
	usedMiB := memory.Used / mib

	return memory.UsedPercent, totalMiB, usedMiB, nil
}
