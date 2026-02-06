package cpu

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessStats struct {
	PID           int32
	Name          string
	CPUPercent    float64
	MemoryPercent float64
	MemoryRSS     uint64
	CreateTime    int64
	Uptime        time.Duration
	Status        string
	Username      string
	Exe           string
	Cmdline       string
}

func GetProcessStats(pid int32) (*ProcessStats, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("процесс с PID %d не найден: %v", pid, err)
	}

	stats := &ProcessStats{
		PID: pid,
	}

	if name, err := proc.Name(); err == nil {
		stats.Name = name
	}

	if cpu, err := proc.CPUPercent(); err == nil {
		stats.CPUPercent = cpu
	}

	if mem, err := proc.MemoryPercent(); err == nil {
		stats.MemoryPercent = float64(mem)
	}

	if memInfo, err := proc.MemoryInfo(); err == nil {
		stats.MemoryRSS = memInfo.RSS
	}

	if createTime, err := proc.CreateTime(); err == nil {
		stats.CreateTime = createTime
		if createTime > 0 {
			stats.Uptime = time.Since(time.Unix(createTime/1000, 0))
		}
	}

	if status, err := proc.Status(); err == nil && len(status) > 0 {
		stats.Status = status[0]
	}

	if username, err := proc.Username(); err == nil {
		stats.Username = username
	}

	if exe, err := proc.Exe(); err == nil {
		stats.Exe = exe
	}

	if cmdline, err := proc.Cmdline(); err == nil {
		stats.Cmdline = cmdline
	}

	return stats, nil
}

func PrintProcessStats(pid int32) {
	stats, err := GetProcessStats(pid)
	if err != nil {
		log.Printf("Ошибка получения статистики процесса: %v", err)
		return
	}

	fmt.Println("\n==========================================")
	fmt.Printf("Статистика процесса PID: %d\n", stats.PID)
	fmt.Println("==========================================")
	fmt.Printf("Имя:              %s\n", stats.Name)
	fmt.Printf("Пользователь:     %s\n", stats.Username)
	fmt.Printf("Статус:           %s\n", stats.Status)
	fmt.Printf("CPU:              %.2f%%\n", stats.CPUPercent)
	fmt.Printf("Память:           %.2f%%\n", stats.MemoryPercent)
	fmt.Printf("Память RSS:       %s\n", formatBytes(stats.MemoryRSS))
	if stats.Uptime > 0 {
		fmt.Printf("Время работы:     %s\n", formatDuration(stats.Uptime))
	}
	if stats.CreateTime > 0 {
		createTime := time.Unix(stats.CreateTime/1000, 0)
		fmt.Printf("Время создания:   %s\n", createTime.Format("2006-01-02 15:04:05"))
	}
	if stats.Exe != "" {
		fmt.Printf("Исполняемый файл: %s\n", stats.Exe)
	}
	if stats.Cmdline != "" {
		fmt.Printf("Командная строка: %s\n", stats.Cmdline)
	}
	fmt.Println("==========================================")
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f сек", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f мин", d.Minutes())
	}
	if d < 24*time.Hour {
		hours := d.Hours()
		minutes := d.Minutes() - float64(int(hours))*60
		return fmt.Sprintf("%.0f ч %.0f мин", hours, minutes)
	}
	days := d.Hours() / 24
	hours := d.Hours() - days*24
	return fmt.Sprintf("%.0f дн %.0f ч", days, hours)
}
