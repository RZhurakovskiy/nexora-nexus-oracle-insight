import type { ApiClient } from "@/services/http/ApiClient"
import type {
	NetworkConnection,
	NetworkInterfaceStat,
	NetworkProcessStat,
} from "@/types/network"

export type ConnectionsKind = "all" | "tcp" | "udp"

export class NetworkService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить список соединений (TCP/UDP).
	 */
	getConnections(kind: ConnectionsKind = "all"): Promise<NetworkConnection[]> {
		const qs = new URLSearchParams()
		qs.set("kind", kind)
		return this.api.getJson<NetworkConnection[]>(
			`/api/network-connections?${qs.toString()}`
		)
	}

	/**
	 * Топ процессов по сетевой активности (по числу соединений/статусов).
	 */
	getTopProcesses(limit = 20): Promise<NetworkProcessStat[]> {
		const qs = new URLSearchParams()
		qs.set("limit", String(limit))
		return this.api.getJson<NetworkProcessStat[]>(
			`/api/network-top-processes?${qs.toString()}`
		)
	}

	/**
	 * Статистика по сетевым интерфейсам.
	 */
	getInterfaces(): Promise<NetworkInterfaceStat[]> {
		return this.api.getJson<NetworkInterfaceStat[]>("/api/network-interfaces")
	}
}
