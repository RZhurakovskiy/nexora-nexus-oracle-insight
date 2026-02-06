import type { ApiClient } from "@/services/http/ApiClient"

type RootStatusResponse = {
	rootStatus: boolean
}

export class RootService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить сводную информацию о системе.
	 */
	getRootService(): Promise<RootStatusResponse> {
		return this.api.getJson<RootStatusResponse>("/api/get-root-status")
	}
}
