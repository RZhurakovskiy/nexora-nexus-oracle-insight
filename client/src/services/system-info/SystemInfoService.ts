import type { ApiClient } from "@/services/http/ApiClient"
import type { SystemInfoResponse } from "@/types/systemInfo"

export class SystemInfoService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить сводную информацию о системе.
	 */
	getSystemInfo(): Promise<SystemInfoResponse> {
		return this.api.getJson<SystemInfoResponse>("/api/system-info")
	}
}
