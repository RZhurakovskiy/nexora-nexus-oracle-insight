import type { WsStatus } from "@/services/realtime/RealtimeSocket"
import type { CpuSample, MemorySample } from "@/types/metrics"
import { defineStore } from "pinia"

const MAX_DATA_POINTS = 60

type MetricsState = {
	/** Текущее состояние мониторинга. */
	monitoringEnabled: boolean
	monitoringMessage: string | null
	monitoringUpdatedAt: string | null

	cpuWsStatus: WsStatus
	memoryWsStatus: WsStatus

	latestCpu: CpuSample | null
	latestMemory: MemorySample | null

	cpuSeries: CpuSample[]
	memorySeries: MemorySample[]

	/** Сообщения для отладки, что ws жив. */
	latestCpuRaw: string | null
	latestMemoryRaw: string | null

	/** Флаг, что держатся активные ws-подключения. */
	socketsConnected: boolean

	/** Ошибка последней операции start/stop/init если была. */
	error: string | null

	starting: boolean
	stopping: boolean
	initializing: boolean
}

export const useMetricsStore = defineStore("metrics", {
	state: (): MetricsState => ({
		monitoringEnabled: false,
		monitoringMessage: null,
		monitoringUpdatedAt: null,

		cpuWsStatus: "idle",
		memoryWsStatus: "idle",

		latestCpu: null,
		latestMemory: null,

		cpuSeries: [],
		memorySeries: [],

		latestCpuRaw: null,
		latestMemoryRaw: null,

		socketsConnected: false,

		error: null,
		starting: false,
		stopping: false,
		initializing: false,
	}),

	getters: {
		/**
		 * Явно показываю "можно ли нажать кнопку старт".
		 */
		canStart(state): boolean {
			return !state.starting && !state.stopping
		},
		canStop(state): boolean {
			return !state.starting && !state.stopping
		},
	},

	actions: {
		/**
		 * Инициализация страницы: читает текущее состояние мониторинга с сервера.
		 */
		async init(): Promise<void> {
			if (this.initializing) return
			this.initializing = true
			this.error = null

			try {
				const res = await this.$services.monitoring.getStatus()
				this.monitoringEnabled = !!res.enabled
				this.monitoringMessage = res.message ?? null
				this.monitoringUpdatedAt = res.timestamp ?? null
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить состояние мониторинга"
			} finally {
				this.initializing = false
			}
		},

		/**
		 * Включить мониторинг:
		 *
		 * enabled = true
		 * запуск ws подключения
		 *
		 */
		async startMonitoring(): Promise<void> {
			if (this.starting || this.stopping) return
			this.starting = true
			this.error = null

			try {
				const res = await this.$services.monitoring.setEnabled(true)
				this.monitoringEnabled = !!res.enabled
				this.monitoringMessage = res.message ?? null
				this.monitoringUpdatedAt = res.timestamp ?? null

				if (this.monitoringEnabled) {
					this.connectSockets()
				}
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось включить мониторинг"
			} finally {
				this.starting = false
			}
		},

		/**
		 * Выключить мониторинг:
		 * enabled = false
		 * отключение ws
		 */
		async stopMonitoring(): Promise<void> {
			if (this.starting || this.stopping) return
			this.stopping = true
			this.error = null

			try {
				this.clearLatestData()

				const res = await this.$services.monitoring.setEnabled(false)
				this.monitoringEnabled = !!res.enabled
				this.monitoringMessage = res.message ?? null
				this.monitoringUpdatedAt = res.timestamp ?? null
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось выключить мониторинг"
			} finally {
				this.stopping = false
			}
		},

		/**
		 * Запуск ws подключения к "/ws/cpu" и "/ws/memory".
		 * Держу это отдельным методом, чтобы можно было отключать WS при уходе со страницы,
		 * не меняя состояние мониторинга на сервере.
		 */
		connectSockets(): void {
			if (this.socketsConnected) return

			this.$services.metricsRealtime.start({
				onCpuStatus: s => (this.cpuWsStatus = s),
				onMemoryStatus: s => (this.memoryWsStatus = s),

				onMonitoringEnabledFromServer: enabled => {
					this.monitoringEnabled = enabled
					if (!enabled) {
						this.clearLatestData()
					}
				},

				onCpuSample: sample => {
					this.latestCpu = sample
					this.pushCpuSample(sample)
				},
				onMemorySample: sample => {
					this.latestMemory = sample
					this.pushMemorySample(sample)
				},

				onSocketMessageError: (scope, raw) => {
					const text = typeof raw === "string" ? raw : JSON.stringify(raw)
					if (scope === "cpu") this.latestCpuRaw = text
					else this.latestMemoryRaw = text
				},
			})

			this.socketsConnected = true
		},

		/**
		 * Отключает ws подключения.
		 */
		disconnectSockets(): void {
			this.$services.metricsRealtime.stop()
			this.socketsConnected = false
			this.cpuWsStatus = "disconnected"
			this.memoryWsStatus = "disconnected"
		},

		/**
		 * Очистить текущие значения CPU/RAM на странице.
		 */
		clearLatestData(): void {
			this.latestCpu = null
			this.latestMemory = null
			this.latestCpuRaw = null
			this.latestMemoryRaw = null
			this.cpuSeries = []
			this.memorySeries = []
		},

		/**
		 * Добавить точку CPU в ring-buffer.
		 */
		pushCpuSample(sample: CpuSample): void {
			const next = [...this.cpuSeries, sample]
			if (next.length > MAX_DATA_POINTS) next.shift()
			this.cpuSeries = next
		},

		/**
		 * Добавить точку RAM в ring-buffer.
		 */
		pushMemorySample(sample: MemorySample): void {
			const next = [...this.memorySeries, sample]
			if (next.length > MAX_DATA_POINTS) next.shift()
			this.memorySeries = next
		},
	},
})
