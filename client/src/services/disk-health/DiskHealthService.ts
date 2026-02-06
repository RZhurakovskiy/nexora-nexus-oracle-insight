import type { ApiClient } from "@/services/http/ApiClient"
import type { DiskHealthResponse } from "@/types/diskHealth"

export class DiskHealthService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить SMART/health информацию по дискам.
	 */
	getDiskHealth(): Promise<DiskHealthResponse> {
		return this.api.getJson<DiskHealthResponse>("/api/disk-health")
	}
}
