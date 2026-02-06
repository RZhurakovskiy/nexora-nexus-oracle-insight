import { getBackendWsBaseUrl } from "@/config/backend"
import type { WsStatus } from "@/services/realtime/RealtimeSocket"
import { RealtimeSocket } from "@/services/realtime/RealtimeSocket"
import type { MonitoringStatusWsPayload, ProcessInfo } from "@/types/processes"

export type ProcessesRealtimeHandlers = {
	onStatus?: (status: WsStatus) => void
	onMonitoringEnabledFromServer?: (enabled: boolean) => void
	onProcesses?: (items: ProcessInfo[]) => void
	onSocketMessageError?: (raw: unknown) => void
}

/**
 * Клиент реального времени для вкладки "Процессы".
 *
 * ws endpoint: "/ws/processes"
 * Отдает:
 * - массив "ProcessInfo[]" - список процессов
 * - либо объект статуса "{ monitoringEnabled: boolean, message?: string }"
 */
export class ProcessesRealtimeClient {
	private socket: RealtimeSocket | null = null
	private readonly wsBaseUrl: string

	constructor(wsBaseUrl: string = getBackendWsBaseUrl()) {
		this.wsBaseUrl = wsBaseUrl
	}

	start(handlers: ProcessesRealtimeHandlers): void {
		this.stop()

		this.socket = new RealtimeSocket(`${this.wsBaseUrl}/ws/processes`)
		this.socket.setHandlers({
			onStatusChange: s => handlers.onStatus?.(s),
			onMessage: evt => this.handleMessage(evt, handlers),
		})
		this.socket.start()
	}

	stop(): void {
		this.socket?.stop()
		this.socket = null
	}

	private handleMessage(
		evt: MessageEvent,
		handlers: ProcessesRealtimeHandlers
	): void {
		const msg = this.safeJsonParse<unknown>(evt.data)
		if (msg == null) {
			return
		}

		if (Array.isArray(msg)) {
			handlers.onProcesses?.(msg as ProcessInfo[])
			return
		}

		if (msg && typeof msg === "object") {
			const maybeStatus = msg as Partial<MonitoringStatusWsPayload>
			if (typeof maybeStatus.monitoringEnabled === "boolean") {
				handlers.onMonitoringEnabledFromServer?.(maybeStatus.monitoringEnabled)
			}
			return
		}

		handlers.onSocketMessageError?.(msg)
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
