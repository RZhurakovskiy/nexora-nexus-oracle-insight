/* Настройка всех HTTP роутов API и WebSocket соединений */
package api

import (
	"net/http"

	"github.com/RZhurakovskiy/agent/server/handlers"
	"github.com/RZhurakovskiy/agent/server/ws"
)

/* Регистрирует все HTTP эндпоинты и WebSocket роуты в мультиплексоре */
func SetupRoutes(mux *http.ServeMux) {

	/* API для получения имени пользователя хоста */
	mux.HandleFunc("/api/gethostusername", handlers.GetHostUserName)
	/* API для завершения процесса по его PID */
	mux.HandleFunc("/api/kill-process-by-id", handlers.KillProcessById)
	/* API для получения имени пользователя хоста */
	mux.HandleFunc("/api/get-host-username", handlers.GetHostUserName)

	/* API для получения детальной информации о процессоре */
	mux.HandleFunc("/api/get-device-info", handlers.GetDeviceInfo)
	/* API для получения общей информации о системе */
	mux.HandleFunc("/api/system-info", handlers.GetSystemInfo)
	/* API для получения информации о здоровье дисков через SMART */
	mux.HandleFunc("/api/disk-health", handlers.GetDiskHealth)
	/* API для получения версии сервера */
	mux.HandleFunc("/api/version", handlers.GetVersion)

	/* API для получения списка портов в состоянии LISTEN */
	mux.HandleFunc("/api/listening-ports", handlers.GetListeningPort)
	/* API для получения всех сетевых соединений */
	mux.HandleFunc("/api/network-connections", handlers.GetNetworkConnections)
	/* API для получения топ процессов по количеству сетевых соединений */
	mux.HandleFunc("/api/network-top-processes", handlers.GetNetworkTopProcesses)
	/* API для получения информации о сетевых интерфейсах */
	mux.HandleFunc("/api/network-interfaces", handlers.GetNetworkInterfaces)

	/* API для запуска нового процесса */
	mux.HandleFunc("/api/start-processes", handlers.StartProcess)

	/* API для экспорта списка процессов в CSV или JSON */
	mux.HandleFunc("/api/export/processes", handlers.ExportProcesses)
	/* API для экспорта текущих метрик CPU и памяти в CSV или JSON */
	mux.HandleFunc("/api/export/metrics", handlers.ExportMetrics)

	/* API для получения истории метрик за указанный период */
	mux.HandleFunc("/api/metrics-history", handlers.GetMetricsHistory)
	/* API для очистки всей истории метрик */
	mux.HandleFunc("/api/clear-metrics", handlers.ClearMetrics)

	/* API для начала записи процессов с высоким использованием ресурсов */
	mux.HandleFunc("/api/start-recording", handlers.StartRecording)
	/* API для остановки записи процессов */
	mux.HandleFunc("/api/stop-recording", handlers.StopRecording)
	/* API для получения статуса текущей сессии записи */
	mux.HandleFunc("/api/recording-status", handlers.GetRecordingStatus)
	/* API для получения списка записанных процессов */
	mux.HandleFunc("/api/recorded-processes", handlers.GetRecordedProcesses)

	/* API для проверки наличия root прав */
	mux.HandleFunc("/api/get-root-status", handlers.GetRootStatus)

	/* API для получения списка алертов, поддерживает только GET метод */
	mux.HandleFunc("/api/alerts", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.GetAlerts(writer, request)
		default:
			http.Error(writer, "Метод не разрешён. Используйте GET", http.StatusMethodNotAllowed)
		}
	})

	/* API для подтверждения алерта */
	mux.HandleFunc("/api/alerts/acknowledge", handlers.AcknowledgeAlert)
	/* API для получения и установки порогов алертов, поддерживает GET и POST методы */
	mux.HandleFunc("/api/alerts/thresholds", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.GetAlertThresholds(writer, request)
		case http.MethodPost:
			handlers.SetAlertThresholds(writer, request)
		default:
			http.Error(writer, "Метод не разрешён. Используйте GET или POST", http.StatusMethodNotAllowed)
		}
	})

	/* API для получения и установки статуса мониторинга, поддерживает GET и POST методы */
	mux.HandleFunc("/api/monitoring-status", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			handlers.GetMonitoringStatus(writer, request)
		case http.MethodPost:
			handlers.SetMonitoringStatus(writer, request)
		default:
			http.Error(writer, "Метод не разрешён. Используйте GET или POST", http.StatusMethodNotAllowed)
		}
	})

	/* WebSocket для потоковой передачи метрик CPU в реальном времени */
	mux.HandleFunc("/ws/cpu", ws.StreamCPU)
	/* WebSocket для потоковой передачи метрик памяти в реальном времени */
	mux.HandleFunc("/ws/memory", ws.StreamMemory)
	/* WebSocket для потоковой передачи списка процессов в реальном времени */
	mux.HandleFunc("/ws/processes", ws.StreamProcesses)
}
