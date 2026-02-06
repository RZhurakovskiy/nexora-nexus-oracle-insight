import type { ApiClient } from "@/services/http/ApiClient"
import type {
	ListeningPort,
	StartProcessRequest,
	StartProcessResponse,
} from "@/types/ports"

export class PortsService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить список TCP соединений/портов с привязкой к PID и процессу.
	 */
	getListeningPorts(): Promise<ListeningPort[]> {
		return this.api.getJson<ListeningPort[]>("/api/listening-ports")
	}

	/**
	 * Запустить новый процесс на вкладке "Сетевые порты".
	 */
	startProcess(req: StartProcessRequest): Promise<StartProcessResponse> {
		if (!req.command || !req.command.trim()) {
			throw new Error("Укажите команду для запуска")
		}
		return this.api.postJson<StartProcessResponse>("/api/start-processes", req)
	}
}
