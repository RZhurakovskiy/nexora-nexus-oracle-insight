/* Обработчики для получения общей информации о системе */
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
)

/* Структура для хранения информации о хосте и операционной системе */
type SystemInfoHost struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platformFamily"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	KernelArch      string `json:"kernelArch"`
	UptimeSec       uint64 `json:"uptimeSec"`
	BootTime        uint64 `json:"bootTime"`
}

/* Структура для хранения информации о процессоре */
type SystemInfoCPU struct {
	ModelName     string `json:"modelName"`
	VendorID      string `json:"vendorId"`
	PhysicalCores int32  `json:"physicalCores"`
	LogicalCores  int32  `json:"logicalCores"`
}

/* Структура для хранения информации об использовании памяти */
type SystemInfoMemory struct {
	TotalBytes      uint64  `json:"totalBytes"`
	AvailableBytes  uint64  `json:"availableBytes"`
	UsedBytes       uint64  `json:"usedBytes"`
	UsedPercent     float64 `json:"usedPercent"`
	SwapTotalBytes  uint64  `json:"swapTotalBytes"`
	SwapUsedBytes   uint64  `json:"swapUsedBytes"`
	SwapUsedPercent float64 `json:"swapUsedPercent"`
}

/* Структура для хранения средней нагрузки системы */
type SystemInfoLoad struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

/* Структура для хранения информации о дисковом разделе */
type SystemInfoDisk struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	TotalBytes  uint64  `json:"totalBytes"`
	UsedBytes   uint64  `json:"usedBytes"`
	FreeBytes   uint64  `json:"freeBytes"`
	UsedPercent float64 `json:"usedPercent"`
}

/* Структура для хранения информации о сетевом интерфейсе */
type SystemInfoNetInterface struct {
	Name  string   `json:"name"`
	Addrs []string `json:"addrs"`
}

/* Структура полного ответа с информацией о системе */
type SystemInfoResponse struct {
	Host       SystemInfoHost           `json:"host"`
	CPU        SystemInfoCPU            `json:"cpu"`
	Memory     SystemInfoMemory         `json:"memory"`
	Load       *SystemInfoLoad          `json:"load"`
	Disks      []SystemInfoDisk         `json:"disks"`
	Interfaces []SystemInfoNetInterface `json:"interfaces"`
}

/* Собирает и возвращает полную информацию о системе включая хост, cpu, память, диски и сетевые интерфейсы */
func GetSystemInfo(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("Ошибка получения host.Info: %v", err)
		http.Error(writer, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Printf("Ошибка получения cpu.Info: %v", err)
	}
	physicalCores, err := cpu.Counts(false)
	if err != nil {
		log.Printf("Ошибка получения количества физических ядер: %v", err)
		physicalCores = 0
	}
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		log.Printf("Ошибка получения количества логических ядер: %v", err)
		logicalCores = 0
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Ошибка получения mem.VirtualMemory: %v", err)
		http.Error(writer, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Ошибка получения mem.SwapMemory: %v", err)
	}

	var loadInfo *SystemInfoLoad
	if avg, err := load.Avg(); err == nil && avg != nil {
		loadInfo = &SystemInfoLoad{Load1: avg.Load1, Load5: avg.Load5, Load15: avg.Load15}
	}

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Ошибка получения disk.Partitions: %v", err)
		partitions = []disk.PartitionStat{}
	}
	disks := make([]SystemInfoDisk, 0, len(partitions))
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil || usage == nil {
			continue
		}
		disks = append(disks, SystemInfoDisk{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	ifaces, err := gnet.Interfaces()
	if err != nil {
		log.Printf("Ошибка получения gnet.Interfaces: %v", err)
		ifaces = []gnet.InterfaceStat{}
	}
	netIfaces := make([]SystemInfoNetInterface, 0, len(ifaces))
	for _, ni := range ifaces {
		addrs := make([]string, 0, len(ni.Addrs))
		for _, a := range ni.Addrs {
			if a.Addr != "" {
				addrs = append(addrs, a.Addr)
			}
		}
		netIfaces = append(netIfaces, SystemInfoNetInterface{Name: ni.Name, Addrs: addrs})
	}

	cpuModel := ""
	cpuVendor := ""
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
		cpuVendor = cpuInfo[0].VendorID
	}

	resp := SystemInfoResponse{
		Host: SystemInfoHost{
			Hostname:        hostInfo.Hostname,
			OS:              hostInfo.OS,
			Platform:        hostInfo.Platform,
			PlatformFamily:  hostInfo.PlatformFamily,
			PlatformVersion: hostInfo.PlatformVersion,
			KernelVersion:   hostInfo.KernelVersion,
			KernelArch:      hostInfo.KernelArch,
			UptimeSec:       hostInfo.Uptime,
			BootTime:        hostInfo.BootTime,
		},
		CPU: SystemInfoCPU{
			ModelName:     cpuModel,
			VendorID:      cpuVendor,
			PhysicalCores: int32(physicalCores),
			LogicalCores:  int32(logicalCores),
		},
		Memory: SystemInfoMemory{
			TotalBytes:      vm.Total,
			AvailableBytes:  vm.Available,
			UsedBytes:       vm.Used,
			UsedPercent:     vm.UsedPercent,
			SwapTotalBytes:  swap.Total,
			SwapUsedBytes:   swap.Used,
			SwapUsedPercent: swap.UsedPercent,
		},
		Load:       loadInfo,
		Disks:      disks,
		Interfaces: netIfaces,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		log.Printf("Ошибка сериализации system-info: %v", err)
	}
}
