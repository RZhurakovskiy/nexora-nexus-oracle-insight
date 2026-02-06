/* Обработчики для получения детальной информации о процессоре */
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
)

/* Структура для хранения детальной информации о процессоре */
type DeviceInfo struct {
	ProcessorName     string   `json:"processor_name"`
	Vendor            string   `json:"vendor"`
	PhysicalCores     int32    `json:"physical_cores"`
	LogicalProcessors int32    `json:"logical_processors"`
	FrequencyMHz      float64  `json:"frequency_mhz"`
	CacheSizeKB       int32    `json:"cache_size_kb"`
	SupportedFlags    []string `json:"supported_flags"`
	Architecture      string   `json:"architecture"`
	Family            string   `json:"family"`
	Model             string   `json:"model"`
}

/*
Собирает детальную информацию о процессоре, архитектуру, флаги и характеристики

	Возвращает структуру с информацией о процессоре или ошибку
*/
func getDeviceInfo() (*DeviceInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	if len(cpuInfo) == 0 {
		return nil, nil
	}

	firstCPU := cpuInfo[0]

	physicalCoreMap := make(map[string]bool)
	for _, cpu := range cpuInfo {
		physicalCoreMap[cpu.CoreID] = true
	}
	physicalCores := int32(len(physicalCoreMap))
	logicalProcessors := int32(len(cpuInfo))

	architecture := "unknown"
	for _, flag := range firstCPU.Flags {
		if flag == "lm" {
			architecture = "x86_64"
			break
		}
		if flag == "pae" {
			architecture = "x86_pae"
		} else if architecture == "unknown" {
			architecture = "x86"
		}
	}

	importantFlags := []string{"sse", "sse2", "sse3", "ssse3", "sse4_1", "sse4_2", "avx", "avx2", "aes", "fma", "mmx", "vmx", "svm"}
	supportedFlags := []string{}
	flagMap := make(map[string]bool)
	for _, flag := range firstCPU.Flags {
		flagMap[flag] = true
	}

	for _, flag := range importantFlags {
		if flagMap[flag] {
			supportedFlags = append(supportedFlags, flag)
		}
	}

	return &DeviceInfo{
		ProcessorName:     firstCPU.ModelName,
		Vendor:            firstCPU.VendorID,
		PhysicalCores:     physicalCores,
		LogicalProcessors: logicalProcessors,
		FrequencyMHz:      firstCPU.Mhz,
		CacheSizeKB:       firstCPU.CacheSize,
		SupportedFlags:    supportedFlags,
		Architecture:      architecture,
		Family:            firstCPU.Family,
		Model:             firstCPU.Model,
	}, nil
}

/* Обрабатывает HTTP запрос и возвращает детальную информацию о процессоре в формате JSON */
func GetDeviceInfo(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	deviceInfo, err := getDeviceInfo()
	if err != nil {
		log.Printf("Ошибка получения информации о процессоре: %v", err)
		http.Error(writer, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if deviceInfo == nil {
		http.Error(writer, "Информация о процессоре не доступна", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(deviceInfo); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
	}
}
