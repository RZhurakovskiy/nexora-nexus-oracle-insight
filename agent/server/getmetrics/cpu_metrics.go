/* Функции для получения метрик использования CPU */
package getmetrics

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

/* Получает процент использования CPU за указанный период времени */
/* Возвращает среднее значение использования CPU или ошибку */
func UsageCPU(duration time.Duration) (float64, error) {

	percents, err := cpu.Percent(duration, false)
	if err != nil || len(percents) == 0 {
		return 0, err
	}
	return percents[0], nil
}
