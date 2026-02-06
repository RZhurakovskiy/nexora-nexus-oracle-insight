import { ApiClient } from "@/services/http/ApiClient"

export type MonitoringStatusResponse = {
	enabled: boolean
	message: string
	timestamp: string
}

export class MonitoringService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить текущее состояние мониторинга.
	 */
	getStatus(): Promise<MonitoringStatusResponse> {
		return this.api.getJson<MonitoringStatusResponse>("/api/monitoring-status")
	}

	/**
	 * Включить или выключить мониторинг.
	 */
	setEnabled(enabled: boolean): Promise<MonitoringStatusResponse> {
		return this.api.postJson<MonitoringStatusResponse>(
			"/api/monitoring-status",
			{ enabled }
		)
	}
}
