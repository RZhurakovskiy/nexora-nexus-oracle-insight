import { getBackendWsBaseUrl } from "@/config/backend"
import type {
	CpuSample,
	CpuWsPayload,
	MemorySample,
	MemoryWsPayload,
} from "@/types/metrics"
import type { WsStatus } from "./RealtimeSocket"
import { RealtimeSocket } from "./RealtimeSocket"

export type MetricsRealtimeHandlers = {
	onCpuStatus?: (status: WsStatus) => void
	onMemoryStatus?: (status: WsStatus) => void
	onMonitoringEnabledFromServer?: (enabled: boolean) => void
	onCpuSample?: (sample: CpuSample) => void
	onMemorySample?: (sample: MemorySample) => void
	onSocketMessageError?: (scope: "cpu" | "memory", raw: unknown) => void
}

/**
 * Клиент реального времени для страницы "Метрики".
 *
 * Держит два сокета:
 * - /ws/cpu
 * - /ws/memory
 */
export class MetricsRealtimeClient {
	private cpuSocket: RealtimeSocket | null = null
	private memorySocket: RealtimeSocket | null = null
	private readonly wsBaseUrl: string

	constructor(wsBaseUrl: string = getBackendWsBaseUrl()) {
		this.wsBaseUrl = wsBaseUrl
	}

	start(handlers: MetricsRealtimeHandlers): void {
		this.stop()

		this.cpuSocket = new RealtimeSocket(`${this.wsBaseUrl}/ws/cpu`)
		this.cpuSocket.setHandlers({
			onStatusChange: s => handlers.onCpuStatus?.(s),
			onMessage: evt => this.handleCpuMessage(evt, handlers),
		})
		this.cpuSocket.start()

		this.memorySocket = new RealtimeSocket(`${this.wsBaseUrl}/ws/memory`)
		this.memorySocket.setHandlers({
			onStatusChange: s => handlers.onMemoryStatus?.(s),
			onMessage: evt => this.handleMemoryMessage(evt, handlers),
		})
		this.memorySocket.start()
	}

	stop(): void {
		this.cpuSocket?.stop()
		this.cpuSocket = null
		this.memorySocket?.stop()
		this.memorySocket = null
	}

	private handleCpuMessage(
		evt: MessageEvent,
		handlers: MetricsRealtimeHandlers
	): void {
		const msg = this.safeJsonParse<CpuWsPayload>(evt.data)
		if (!msg) {
			handlers.onSocketMessageError?.("cpu", evt.data)
			return
		}

		if (typeof msg.monitoringEnabled === "boolean") {
			handlers.onMonitoringEnabledFromServer?.(msg.monitoringEnabled)

			if (!msg.monitoringEnabled) return
		}

		if (
			typeof msg.cpu === "number" &&
			typeof msg.timestamp === "string" &&
			msg.timestamp
		) {
			handlers.onCpuSample?.({ cpu: msg.cpu, timestamp: msg.timestamp })
		}
	}

	private handleMemoryMessage(
		evt: MessageEvent,
		handlers: MetricsRealtimeHandlers
	): void {
		const msg = this.safeJsonParse<MemoryWsPayload>(evt.data)
		if (!msg) {
			handlers.onSocketMessageError?.("memory", evt.data)
			return
		}

		if (typeof msg.monitoringEnabled === "boolean") {
			handlers.onMonitoringEnabledFromServer?.(msg.monitoringEnabled)
			if (!msg.monitoringEnabled) return
		}

		if (
			typeof msg.memoryUsage === "number" &&
			typeof msg.usedMB === "number" &&
			typeof msg.totalMemory === "number" &&
			typeof msg.timestamp === "string" &&
			msg.timestamp
		) {
			handlers.onMemorySample?.({
				memoryUsage: msg.memoryUsage,
				usedMB: msg.usedMB,
				totalMemoryMB: msg.totalMemory,
				timestamp: msg.timestamp,
			})
		}
	}

	private safeJsonParse<T>(raw: unknown): T | null {
		if (typeof raw !== "string") return null
		try {
			return JSON.parse(raw) as T
		} catch {
			return null
		}
	}
}
