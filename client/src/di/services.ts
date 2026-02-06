import { getBackendBaseUrl, getBackendWsBaseUrl } from "@/config/backend"
import { AuthService } from "@/services/auth/AuthService"
import { DiskHealthService } from "@/services/disk-health/DiskHealthService"
import { HostInfoService } from "@/services/host-info/HostInfoService"
import { ApiClient } from "@/services/http/ApiClient"
import { MetricsHistoryService } from "@/services/metrics-history/MetricsHistoryService"
import { MonitoringService } from "@/services/monitoring/MonitoringService"
import { NetworkService } from "@/services/network/NetworkService"
import { PortsService } from "@/services/ports/PortsService"
import { ProcessesService } from "@/services/processes/ProcessesService"
import { MetricsRealtimeClient } from "@/services/realtime/MetricsRealtimeClient"
import { ProcessesRealtimeClient } from "@/services/realtime/ProcessesRealtimeClient"
import { RecordingService } from "@/services/recording/RecordingService"
import { SystemInfoService } from "@/services/system-info/SystemInfoService"
import { LoadingService } from "@/services/ui/LoadingService"
import { VersionService } from "@/services/version/VersionService"
import type { AxiosInstance } from "axios"
import axios from "axios"

import { RootService } from "@/services/root-status/RootStatusService"

export type AppServices = {
	axios: AxiosInstance
	api: ApiClient
	loading: LoadingService
	auth: AuthService
	version: VersionService
	monitoring: MonitoringService
	metricsHistory: MetricsHistoryService
	recording: RecordingService
	ports: PortsService
	metricsRealtime: MetricsRealtimeClient
	processesRealtime: ProcessesRealtimeClient
	processes: ProcessesService
	hostInfo: HostInfoService
	systemInfo: SystemInfoService
	network: NetworkService
	diskHealth: DiskHealthService
	rootService: RootService
}

/**
 * Фабрика DI-контейнера сервисов приложения.
 *
 */
export function createAppServices(): AppServices {
	const http = axios.create({
		baseURL: getBackendBaseUrl(),
		timeout: 15_000,
	})

	const api = new ApiClient(http)
	const loading = new LoadingService()

	http.interceptors.request.use(
		config => {
			loading.begin()
			return config
		},
		error => {
			loading.end()
			return Promise.reject(error)
		}
	)

	http.interceptors.response.use(
		response => {
			loading.end()
			return response
		},
		error => {
			loading.end()
			return Promise.reject(error)
		}
	)

	return {
		axios: http,
		api,
		loading,
		auth: new AuthService(),
		version: new VersionService(api),
		monitoring: new MonitoringService(api),
		metricsHistory: new MetricsHistoryService(api),
		recording: new RecordingService(api),
		ports: new PortsService(api),
		metricsRealtime: new MetricsRealtimeClient(getBackendWsBaseUrl()),
		processesRealtime: new ProcessesRealtimeClient(getBackendWsBaseUrl()),
		processes: new ProcessesService(api, http),
		hostInfo: new HostInfoService(api),
		systemInfo: new SystemInfoService(api),
		network: new NetworkService(api),
		diskHealth: new DiskHealthService(api),
		rootService: new RootService(api),
	}
}
