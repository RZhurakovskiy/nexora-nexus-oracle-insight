/* Функции для получения информации о здоровье дисков через SMART */
package getmetrics

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

/* Структура для ответа с информацией о здоровье дисков */
type DiskHealthResponse struct {
	Supported bool               `json:"supported"`
	Message   string             `json:"message"`
	Devices   []DiskHealthDevice `json:"devices"`
}

/* Структура для хранения информации о конкретном диске и его SMART статусе */
type DiskHealthDevice struct {
	Device          string   `json:"device"`
	Model           string   `json:"model"`
	Serial          string   `json:"serial"`
	BusType         string   `json:"busType"`
	SmartAvailable  bool     `json:"smartAvailable"`
	SmartEnabled    bool     `json:"smartEnabled"`
	SmartPassed     *bool    `json:"smartPassed"`
	TemperatureC    *int     `json:"temperatureC"`
	PowerOnHours    *int     `json:"powerOnHours"`
	UnsafeShutdowns *int     `json:"unsafeShutdowns"`
	NVMePercentUsed *int     `json:"nvmePercentUsed"`
	Warnings        []string `json:"warnings"`
}

var (
	reDevSda1  = regexp.MustCompile(`^/dev/sd[a-z]+\d+$`)
	reDevNvmeP = regexp.MustCompile(`^/dev/nvme\d+n\d+p\d+$`)
)

/* Извлекает базовое имя блочного устройства из имени раздела */
/* Убирает номер раздела для SATA дисков и номер партиции для NVMe */
func baseBlockDevice(dev string) string {
	d := strings.TrimSpace(dev)
	if d == "" {
		return ""
	}
	if reDevSda1.MatchString(d) {
		return regexp.MustCompile(`\d+$`).ReplaceAllString(d, "")
	}
	if reDevNvmeP.MatchString(d) {
		return regexp.MustCompile(`p\d+$`).ReplaceAllString(d, "")
	}
	return d
}

/* Получает информацию о здоровье дисков через smartctl для Linux систем */
/* Возвращает структуру с информацией о поддержке SMART и списком устройств с их статусами */
func GetDiskHealth(goos string) DiskHealthResponse {
	if goos != "linux" {
		return DiskHealthResponse{
			Supported: false,
			Message:   "SMART/Health доступен только на Linux (best-effort).",
			Devices:   []DiskHealthDevice{},
		}
	}

	path, err := exec.LookPath("smartctl")
	if err != nil || path == "" {
		return DiskHealthResponse{
			Supported: false,
			Message:   "Не найден smartctl. Установите пакет smartmontools (например: `sudo pacman -S smartmontools`).",
			Devices:   []DiskHealthDevice{},
		}
	}

	parts, _ := disk.Partitions(false)
	set := map[string]struct{}{}
	for _, p := range parts {
		if p.Device == "" {
			continue
		}
		if !strings.HasPrefix(p.Device, "/dev/") {
			continue
		}
		base := baseBlockDevice(p.Device)
		if base == "" || !strings.HasPrefix(base, "/dev/") {
			continue
		}
		set[base] = struct{}{}
	}

	devices := make([]string, 0, len(set))
	for d := range set {
		devices = append(devices, d)
	}
	sort.Strings(devices)

	out := make([]DiskHealthDevice, 0, len(devices))
	for _, dev := range devices {
		if !strings.HasPrefix(dev, "/dev/") {
			continue
		}
		item := probeSmartctl(dev)
		out = append(out, item)
	}

	return DiskHealthResponse{
		Supported: true,
		Message:   "Данные SMART получены через smartctl (best-effort). Для некоторых дисков могут требоваться права root.",
		Devices:   out,
	}
}

type smartctlJSON struct {
	Device struct {
		Name    string `json:"name"`
		Model   string `json:"model_name"`
		Serial  string `json:"serial_number"`
		BusType string `json:"bus_type"`
	} `json:"device"`

	SmartStatus *struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`

	Smartctl *struct {
		ExitStatus int `json:"exit_status"`
	} `json:"smartctl"`

	Temperature *struct {
		Current int `json:"current"`
	} `json:"temperature"`

	NVMeSmartHealth *struct {
		Temperature  int `json:"temperature"`
		PercentUsed  int `json:"percentage_used"`
		UnsafeShut   int `json:"unsafe_shutdowns"`
		PowerOnHours int `json:"power_on_hours"`
	} `json:"nvme_smart_health_information_log"`

	PowerOnTime *struct {
		Hours int `json:"hours"`
	} `json:"power_on_time"`
}

func probeSmartctl(dev string) DiskHealthDevice {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "smartctl", "-j", "-a", dev)
	b, err := cmd.Output()

	item := DiskHealthDevice{
		Device:   dev,
		Model:    "",
		Serial:   "",
		BusType:  "",
		Warnings: []string{},
	}

	if err != nil {
		item.SmartAvailable = false
		item.SmartEnabled = false
		item.Warnings = append(item.Warnings, "SMART недоступен (возможно нужны права root или устройство не поддерживается)")
		return item
	}

	var doc smartctlJSON
	if jsonErr := json.Unmarshal(b, &doc); jsonErr != nil {
		item.SmartAvailable = false
		item.SmartEnabled = false
		item.Warnings = append(item.Warnings, "Не удалось распарсить JSON от smartctl")
		return item
	}

	item.Model = doc.Device.Model
	item.Serial = doc.Device.Serial
	item.BusType = doc.Device.BusType
	item.SmartAvailable = true
	item.SmartEnabled = true

	if doc.SmartStatus != nil {
		v := doc.SmartStatus.Passed
		item.SmartPassed = &v
		if !v {
			item.Warnings = append(item.Warnings, "SMART-статус: неисправность (FAILED)")
		}
	}

	if doc.NVMeSmartHealth != nil {
		t := doc.NVMeSmartHealth.Temperature
		item.TemperatureC = &t
		pu := doc.NVMeSmartHealth.PercentUsed
		item.NVMePercentUsed = &pu
		ush := doc.NVMeSmartHealth.UnsafeShut
		item.UnsafeShutdowns = &ush
		poh := doc.NVMeSmartHealth.PowerOnHours
		item.PowerOnHours = &poh
		if pu >= 95 {
			item.Warnings = append(item.Warnings, "Износ NVMe: ≥ 95%")
		}
	} else {
		if doc.Temperature != nil && doc.Temperature.Current != 0 {
			t := doc.Temperature.Current
			item.TemperatureC = &t
		}
		if doc.PowerOnTime != nil && doc.PowerOnTime.Hours > 0 {
			h := doc.PowerOnTime.Hours
			item.PowerOnHours = &h
		}
	}

	if doc.Smartctl != nil && doc.Smartctl.ExitStatus != 0 {
		item.Warnings = append(item.Warnings, "smartctl: exit_status != 0 (есть предупреждения)")
	}

	if doc.Device.Name != "" && filepath.Clean(doc.Device.Name) != filepath.Clean(dev) {
		item.Warnings = append(item.Warnings, "Имя устройства в SMART отличается: "+doc.Device.Name)
	}

	return item
}
