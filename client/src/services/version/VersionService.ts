import type { ApiClient } from "@/services/http/ApiClient"
import type { VersionResponse } from "@/types/version"

/**
 * Сервис для получения версии backend/агента.
 */
export class VersionService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Возвращает версию backend.
	 */
	async getVersion(): Promise<VersionResponse> {
		return this.api.getJson<VersionResponse>("/api/version")
	}
}
