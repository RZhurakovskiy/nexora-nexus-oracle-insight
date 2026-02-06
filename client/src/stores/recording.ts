import { useNotificationsStore } from "@/stores/notifications"
import type { RecordingSession, StartRecordingParams } from "@/types/recording"
import { defineStore } from "pinia"

type RecordingState = {
	active: boolean
	session: RecordingSession | null
	remainingSeconds: number | null
	starting: boolean
	stopping: boolean
	error: string | null
	notifiedSessionEndId: number | null
	lastSessionId: number | null
}

let tickTimer: number | null = null

export const useRecordingStore = defineStore("recording", {
	state: (): RecordingState => ({
		active: false,
		session: null,
		remainingSeconds: null,
		starting: false,
		stopping: false,
		error: null,
		notifiedSessionEndId: null,
		lastSessionId: null,
	}),

	getters: {
		/**
		 * Форматирование оставшегося времени как "MM:SS".
		 */
		remainingText(state): string {
			if (state.remainingSeconds == null) return "—"
			const s = Math.max(0, Math.floor(state.remainingSeconds))
			const mm = Math.floor(s / 60)
			const ss = s % 60
			return `${String(mm).padStart(2, "0")}:${String(ss).padStart(2, "0")}`
		},
	},

	actions: {
		/**
		 * Парс даты, которые backend отдаёт строкой `YYYY-MM-DD HH:mm:ss`.
		 */
		parseServerDateTime(value: string): Date | null {
			const v = String(value ?? "").trim()
			if (!v) return null

			let d = new Date(v)
			if (!Number.isNaN(d.getTime())) return d

			const isoLike = v.replace(" ", "T")
			d = new Date(isoLike)
			return Number.isNaN(d.getTime()) ? null : d
		},

		/**
		 * Пересчитать оставшееся время по `session.endTime`.
		 */
		recomputeRemaining(): void {
			if (!this.session) {
				this.remainingSeconds = null
				return
			}

			const end = this.parseServerDateTime(this.session.endTime)
			if (!end) {
				this.remainingSeconds = null
				return
			}

			const diffMs = end.getTime() - Date.now()
			this.remainingSeconds = Math.max(0, Math.ceil(diffMs / 1000))
		},

		/**
		 * Разово получить статус записи с сервера.
		 */
		async fetchStatus(): Promise<void> {
			this.error = null
			try {
				const res = await this.$services.recording.getStatus()
				this.active = !!res.active
				this.session = res.session ?? null
				this.lastSessionId = this.session?.id ?? this.lastSessionId
				this.recomputeRemaining()
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить статус записи"
			}
		},

		/**
		 * Запустить локальный тик таймера по `endTime`.
		 *
		 * Мы не дергаем `recording-status` интервалом: backend не пушит события,
		 * а для UX достаточно локального countdown + 1 подтверждения статуса в момент окончания.
		 */
		startTick(): void {
			if (tickTimer != null) window.clearInterval(tickTimer)
			tickTimer = window.setInterval(() => {
				this.recomputeRemaining()
				this.maybeNotifySessionEnded()
			}, 1000)
		},

		/**
		 * Остановить таймер.
		 */
		stopTick(): void {
			if (tickTimer != null) window.clearInterval(tickTimer)
			tickTimer = null
		},

		/**
		 * Инициализация: читаем текущий статус и запускаем локальный таймер.
		 */
		async init(): Promise<void> {
			await this.fetchStatus()
			this.startTick()
		},

		/**
		 * Показать уведомление, когда сессия завершилась по таймеру.
		 */
		async maybeNotifySessionEnded(): Promise<void> {
			if (!this.active) return
			if (this.stopping || this.starting) return
			if (!this.session) return
			if (this.remainingSeconds == null) return
			if (this.remainingSeconds > 0) return

			const id = this.session.id
			if (this.notifiedSessionEndId === id) return

			/**
			 * Подтверждаем статус одним запросом, чтобы не показывать уведомление раньше, чем backend закрыл сессию.
			 */
			await this.fetchStatus()
			if (this.active) return

			this.notifiedSessionEndId = id
			useNotificationsStore().success(
				`Сессия #${id} завершена. Записи доступны во вкладке «История метрик».`,
				"Запись метрик"
			)
		},

		/**
		 * Запуск записи метрик.
		 *
		 * Важно: валидирую `duration >= 60` на фронте, потому что backend сейчас принимает любое >0.
		 */
		async start(params: StartRecordingParams): Promise<void> {
			if (this.starting || this.stopping) return
			if (this.active) {
				useNotificationsStore().info("Запись уже активна", "Запись метрик")
				return
			}

			this.starting = true
			this.error = null

			try {
				if (params.duration < 60) {
					throw new Error("Минимум 60 секунд")
				}
				const res = await this.$services.recording.start(params)
				useNotificationsStore().success(
					res.message ?? "Запись метрик запущена",
					"Запись метрик"
				)
				await this.fetchStatus()
			} catch (e: any) {
				const msg = e?.message ?? "Не удалось запустить запись метрик"
				this.error = msg
				useNotificationsStore().error(msg, "Запись метрик")
			} finally {
				this.starting = false
			}
		},

		/**
		 * Остановить запись метрик.
		 */
		async stop(): Promise<void> {
			if (this.stopping || this.starting) return
			this.stopping = true
			this.error = null

			try {
				const res = await this.$services.recording.stop()
				useNotificationsStore().success(
					res.message ?? "Запись метрик остановлена",
					"Запись метрик"
				)
				await this.fetchStatus()
			} catch (e: any) {
				const msg = e?.message ?? "Не удалось остановить запись метрик"
				this.error = msg
				useNotificationsStore().error(msg, "Запись метрик")
			} finally {
				this.stopping = false
			}
		},
	},
})
