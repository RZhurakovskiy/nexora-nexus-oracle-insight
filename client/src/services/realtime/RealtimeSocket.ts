export type WsStatus =
	| "idle"
	| "connecting"
	| "connected"
	| "reconnecting"
	| "disconnected"
	| "error"

export type RealtimeSocketHandlers = {
	onOpen?: () => void
	onMessage?: (evt: MessageEvent) => void
	onError?: (evt: Event) => void
	onClose?: (evt: CloseEvent) => void
	onStatusChange?: (status: WsStatus) => void
}

export type RealtimeSocketOptions = {
	maxDelayMs?: number
}

/**
 * Обёртка над ws с автопереподключением.
 */
export class RealtimeSocket {
	private readonly url: string
	private readonly maxDelayMs: number
	private ws: WebSocket | null = null
	private reconnectAttempts = 0
	private reconnectTimer: number | null = null
	private isStopped = true
	private handlers: RealtimeSocketHandlers = {}

	constructor(
		url: string,
		{ maxDelayMs = 10_000 }: RealtimeSocketOptions = {}
	) {
		this.url = url
		this.maxDelayMs = maxDelayMs
	}

	setHandlers(handlers: RealtimeSocketHandlers = {}): void {
		this.handlers = handlers
	}

	start(): void {
		this.isStopped = false
		this.connect()
	}

	stop(): void {
		this.isStopped = true

		if (this.reconnectTimer != null) {
			window.clearTimeout(this.reconnectTimer)
			this.reconnectTimer = null
		}

		if (this.ws) {
			this.ws.onclose = null
			this.ws.onerror = null
			this.ws.onmessage = null
			this.ws.onopen = null

			if (
				this.ws.readyState === WebSocket.OPEN ||
				this.ws.readyState === WebSocket.CONNECTING
			) {
				this.ws.close()
			}
		}

		this.ws = null
		this.setStatus("disconnected")
	}

	send(data: string): void {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(data)
		}
	}

	get readyState(): number | null {
		return this.ws?.readyState ?? null
	}

	private connect(): void {
		if (this.isStopped) return
		try {
			this.setStatus("connecting")
			this.ws = new WebSocket(this.url)
		} catch {
			this.scheduleReconnect()
			return
		}

		this.ws.onopen = () => {
			this.setStatus("connected")
			this.reconnectAttempts = 0
			this.handlers.onOpen?.()
		}

		this.ws.onmessage = evt => this.handlers.onMessage?.(evt)

		this.ws.onerror = evt => {
			this.setStatus("error")
			this.handlers.onError?.(evt)
		}

		this.ws.onclose = evt => {
			this.handlers.onClose?.(evt)
			this.setStatus(this.isStopped ? "disconnected" : "reconnecting")
			if (!this.isStopped) this.scheduleReconnect()
		}
	}

	private scheduleReconnect(): void {
		if (this.isStopped) return
		this.setStatus("reconnecting")
		const delay = Math.min(1000 * 2 ** this.reconnectAttempts, this.maxDelayMs)
		this.reconnectAttempts += 1
		this.reconnectTimer = window.setTimeout(() => this.connect(), delay)
	}

	private setStatus(status: WsStatus): void {
		this.handlers.onStatusChange?.(status)
	}
}
