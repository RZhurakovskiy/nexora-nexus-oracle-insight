import { useNotificationsStore } from "@/stores/notifications"
import type {
	ListeningPort,
	StartProcessRequest,
	StartProcessResponse,
} from "@/types/ports"
import { defineStore } from "pinia"

type PortsState = {
	ports: ListeningPort[]
	loading: boolean
	error: string | null
	initialized: boolean
	starting: boolean
	lastStartResult: StartProcessResponse | null
	killingPid: number | null
}

export const usePortsStore = defineStore("ports", {
	state: (): PortsState => ({
		ports: [],
		loading: false,
		error: null,
		initialized: false,
		starting: false,
		lastStartResult: null,
		killingPid: null,
	}),

	actions: {
		/**
		 * Инициализация страницы.
		 *
		 * "if (this.initialized || this.loading) return" Блокирую выполнение, если загрузка уже идет,
		 * чтобы избежать повторных запросов.
		 */
		async init(): Promise<void> {
			if (this.initialized || this.loading) return
			await this.refresh()
			this.initialized = !this.error
		},

		/**
		 * Принудительно обновить список портов.
		 */
		async refresh(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null

			try {
				const data = await this.$services.ports.getListeningPorts()
				this.ports = Array.isArray(data) ? data : []
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось загрузить список портов"
			} finally {
				this.loading = false
			}
		},

		/**
		 * Запустить процесс и затем обновить список портов.
		 *
		 * Небольшая задержка перед обновлением, потому что новый процесс
		 * может не сразу появиться в списке сетевых соединений.
		 * Но не надежный способ полагаться на setTimeout. Возможно стоит пересмотреть
		 */
		async startProcess(
			req: StartProcessRequest
		): Promise<{ ok: boolean; message: string }> {
			if (this.starting)
				return { ok: false, message: "Операция уже выполняется" }
			this.starting = true
			this.error = null
			this.lastStartResult = null

			try {
				const res = await this.$services.ports.startProcess(req)
				this.lastStartResult = res

				await new Promise(resolve => setTimeout(resolve, 5000))
				await this.refresh()

				return { ok: true, message: res.msg || "Процесс запущен" }
			} catch (e: any) {
				const message = e?.message ?? "Не удалось запустить процесс"
				this.error = message
				return { ok: false, message }
			} finally {
				this.starting = false
			}
		},

		/**
		 * Завершить процесс по PID и обновить список портов.
		 *
		 */
		async killProcess(pid: number): Promise<{ ok: boolean; message: string }> {
			if (!Number.isFinite(pid) || pid <= 0)
				return { ok: false, message: "Некорректный PID" }
			if (this.killingPid)
				return { ok: false, message: "Операция уже выполняется" }

			this.killingPid = pid
			this.error = null

			try {
				const res = await this.$services.processes.killProcessById(pid)

				await new Promise(resolve => setTimeout(resolve, 500))
				await this.refresh()
				useNotificationsStore().success(
					`PID: ${pid}\n${res.message ?? "Процесс завершён"}`,
					"Процесс завершён"
				)
				return { ok: true, message: res.message ?? "Процесс завершён" }
			} catch (e: any) {
				const message = e?.message ?? "Не удалось завершить процесс"
				this.error = message
				useNotificationsStore().error(
					`PID: ${pid}\n${message}`,
					"Не удалось завершить процесс"
				)
				return { ok: false, message }
			} finally {
				this.killingPid = null
			}
		},
	},
})
