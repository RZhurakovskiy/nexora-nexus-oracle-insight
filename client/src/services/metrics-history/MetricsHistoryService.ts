import type { ApiClient } from "@/services/http/ApiClient"
import type { MetricsHistoryItem } from "@/types/metricsHistory"

export type MetricsHistoryQuery = {
	from?: string
	to?: string
	limit?: number
}

export class MetricsHistoryService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить историю метрик.
	 *
	 * "from/to" передается без UTC-конвертации, в формате как в БД "YYYY-MM-DD HH:mm:ss",
	 * иначе фильтр по времени будет некорректный.
	 */
	getHistory(query: MetricsHistoryQuery = {}): Promise<MetricsHistoryItem[]> {
		const params = new URLSearchParams()
		if (query.from) params.set("from", query.from)
		if (query.to) params.set("to", query.to)
		if (typeof query.limit === "number" && query.limit > 0)
			params.set("limit", String(query.limit))
		const qs = params.toString()
		return this.api.getJson<MetricsHistoryItem[]>(
			`/api/metrics-history${qs ? `?${qs}` : ""}`
		)
	}

	/**
	 * Очистить историю метрик.
	 */
	clearHistory(): Promise<{ success: boolean; message: string }> {
		return this.api.postJson<{ success: boolean; message: string }>(
			"/api/clear-metrics",
			{}
		)
	}
}
