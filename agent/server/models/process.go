/* Модели данных для процессов и метрик системы */
package models

/* Структура для хранения информации о системном процессе */
type ProcessInfo struct {
	PID           int32    `json:"pid"`
	Name          string   `json:"name"`
	Exe           string   `json:"exe"`
	Cmdline       string   `json:"cmdline"`
	Username      string   `json:"username"`
	Status        string   `json:"status"`
	CreateTime    int64    `json:"createTime"`
	ParentPID     int32    `json:"parentPid"`
	CPUPercent    float64  `json:"cpuPercent"`
	MemoryPercent float64  `json:"memoryPercent"`
	MemoryRSS     uint64   `json:"memoryRss"`
	Ports         []uint32 `json:"ports"`
}

/* Структура для ответа с метриками использования CPU */
type CPUMetricsResponse struct {
	CPU       float64 `json:"cpu"`
	Timestamp string  `json:"timestamp"`
}

/* Структура для ответа с метриками использования памяти */
type MemoryMetricsResponse struct {
	MemoryUsage float64 `json:"memory"`
	UsedMB      uint64  `json:"usedMB"`
	TotalMemory uint64  `json:"totalmemory"`
	Timestamp   string  `json:"timestamp"`
}

type KillProcessByID struct {
	PID       int32  `json:"pid"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

/* Структура для запроса на изменение состояния мониторинга */
type MonitoringStatusRequest struct {
	Enabled bool `json:"enabled"`
}

/* Структура для ответа о состоянии мониторинга */
type MonitoringStatusResponse struct {
	Enabled   bool   `json:"enabled"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

/* Структура для хранения информации о порте в состоянии LISTEN */
type ListeningPort struct {
	Port       uint32 `json:"port"`
	Protocol   string `json:"protocol"`
	PID        int32  `json:"pid"`
	Process    string `json:"process"`
	Status     string `json:"status"`
	LocalAddr  string `json:"localAddr"`
	RemoteAddr string `json:"remoteAddr"`
}

/* Структура для запроса на запуск нового процесса */
type StartProcessRequest struct {
	Command   string `json:"command"`
	Args      string `json:"args"`
	Cwd       string `json:"cwd"`
	Timestamp string `json:"timestamp"`
}

/* Структура для ответа с результатом запуска процесса */
type StartProcessResponse struct {
	PID     int32  `json:"pid"`
	Command string `json:"command"`
	Args    string `json:"args"`
	Cwd     string `json:"cwd"`
	Msg     string `json:"msg"`
}
