import type { ApiClient } from "@/services/http/ApiClient"
import type {
	RecordingStatusResponse,
	StartRecordingParams,
	StartRecordingResponse,
} from "@/types/recording"

export class RecordingService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить текущий статус записи метрик.
	 */
	getStatus(): Promise<RecordingStatusResponse> {
		return this.api.getJson<RecordingStatusResponse>("/api/recording-status")
	}

	/**
	 * Запустить запись метрик.
	 */
	start(params: StartRecordingParams): Promise<StartRecordingResponse> {
		return this.api.postJson<StartRecordingResponse>(
			"/api/start-recording",
			params
		)
	}

	/**
	 * Остановить запись метрик.
	 */
	stop(): Promise<{ success: boolean; message: string }> {
		return this.api.postJson<{ success: boolean; message: string }>(
			"/api/stop-recording",
			{}
		)
	}
}
