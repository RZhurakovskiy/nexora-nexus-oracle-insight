import type { ExportFormat } from "@/services/processes/ProcessesService"
import type { WsStatus } from "@/services/realtime/RealtimeSocket"
import { useNotificationsStore } from "@/stores/notifications"
import type { ProcessInfo } from "@/types/processes"
import { defineStore } from "pinia"

type ProcessesState = {
	processes: ProcessInfo[]
	wsStatus: WsStatus
	socketsConnected: boolean
	isInitialLoading: boolean
	hasReceivedFirstPayload: boolean
	error: string | null
	killing: boolean
	exporting: boolean
}

export const useProcessesStore = defineStore("processes", {
	state: (): ProcessesState => ({
		processes: [],
		wsStatus: "disconnected",
		socketsConnected: false,
		isInitialLoading: true,
		hasReceivedFirstPayload: false,
		error: null,
		killing: false,
		exporting: false,
	}),

	actions: {
		/**
		 * Перевести UI в режим ожидания первого списка процессов.
		 *
		 * Я использую это при локальном переключении мониторинга, сервер шлёт список раз в 5 сек,
		 * и без этого флага UI будет показывать "Процесс не найден" до получения первых данных.
		 */
		setWaitingForFirstPayload(): void {
			this.isInitialLoading = true
			this.hasReceivedFirstPayload = false
		},

		/**
		 * Запуск WS "/ws/processes" для получения списка процессов.
		 *
		 * Не включает мониторинг на сервере, а только подключает transport.
		 * Включение/выключение мониторинга остаётся за кнопкой в header.
		 */
		connectSockets(): void {
			if (this.socketsConnected) return
			this.error = null
			this.setWaitingForFirstPayload()

			this.$services.processesRealtime.start({
				onStatus: s => (this.wsStatus = s),
				onMonitoringEnabledFromServer: enabled => {
					if (!enabled) {
						this.clearProcesses()
					}
				},
				onProcesses: items => {
					this.processes = Array.isArray(items) ? items : []
					this.hasReceivedFirstPayload = true
					this.isInitialLoading = false
				},
				onSocketMessageError: raw => {
					const text = typeof raw === "string" ? raw : JSON.stringify(raw)
					this.error = `Некорректное сообщение WS: ${text}`
					this.isInitialLoading = false
				},
			})

			this.socketsConnected = true
		},

		/**
		 * Отключить WS подключения.
		 */
		disconnectSockets(): void {
			this.$services.processesRealtime.stop()
			this.socketsConnected = false
			this.wsStatus = "disconnected"
		},

		clearProcesses(): void {
			this.processes = []
			this.hasReceivedFirstPayload = true
			this.isInitialLoading = false
		},

		/**
		 * Завершить процесс по PID.
		 */
		async killProcess(pid: number): Promise<{ ok: boolean; message: string }> {
			if (this.killing)
				return { ok: false, message: "Операция уже выполняется" }
			this.killing = true
			this.error = null

			try {
				const res = await this.$services.processes.killProcessById(pid)
				useNotificationsStore().success(
					`PID: ${pid}\n${res.message ?? "Процесс завершён"}`,
					"Процесс завершён"
				)
				return { ok: true, message: res.message }
			} catch (e: any) {
				const message = e?.message ?? "Не удалось завершить процесс"
				this.error = message
				useNotificationsStore().error(
					`PID: ${pid}\n${message}`,
					"Не удалось завершить процесс"
				)
				return { ok: false, message }
			} finally {
				this.killing = false
			}
		},

		/**
		 * Экспортировать процессы.
		 */
		async exportProcesses(
			format: ExportFormat
		): Promise<{ ok: boolean; message: string }> {
			if (this.exporting)
				return { ok: false, message: "Экспорт уже выполняется" }
			this.exporting = true
			this.error = null

			try {
				await this.$services.processes.exportProcesses(format)
				return { ok: true, message: `Экспортировано в ${format.toUpperCase()}` }
			} catch (e: any) {
				const message = e?.message ?? "Ошибка экспорта процессов"
				this.error = message
				return { ok: false, message }
			} finally {
				this.exporting = false
			}
		},
	},
})
