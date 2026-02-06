import { useNotificationsStore } from "@/stores/notifications"
import type { MetricsHistoryItem } from "@/types/metricsHistory"
import { defineStore } from "pinia"

type MetricsHistoryState = {
	items: MetricsHistoryItem[]
	loading: boolean
	error: string | null

	fromDate: string
	toDate: string
	limit: number
}

export const useMetricsHistoryStore = defineStore("metricsHistory", {
	state: (): MetricsHistoryState => ({
		items: [],
		loading: false,
		error: null,
		fromDate: "",
		toDate: "",
		limit: 1000,
	}),

	actions: {
		/**
		 * Скачать текст как файл.
		 */
		downloadFile(filename: string, mime: string, content: string): void {
			const blob = new Blob([content], { type: mime })
			const url = URL.createObjectURL(blob)
			const a = document.createElement("a")
			a.href = url
			a.download = filename
			document.body.appendChild(a)
			a.click()
			document.body.removeChild(a)
			URL.revokeObjectURL(url)
		},

		/**
		 * Выгрузить записи из БД по текущим фильтрам.
		 * Если записей нет - показать уведомление.
		 */
		async fetchFromDb(): Promise<void> {
			await this.load()
			if (this.items.length === 0) {
				useNotificationsStore().info("Записей не найдено", "История метрик")
			} else {
				useNotificationsStore().success(
					`Загружено записей: ${this.items.length}`,
					"История метрик"
				)
			}
		},

		/**
		 * Экспортировать текущие загруженные записи в JSON/CSV.
		 */
		exportLoaded(format: "json" | "csv" = "json"): void {
			if (!this.items || this.items.length === 0) {
				useNotificationsStore().info("Записей не найдено", "История метрик")
				return
			}

			const ts = new Date().toISOString().slice(0, 19).replace(/:/g, "-")
			if (format === "json") {
				this.downloadFile(
					`metrics_history_${ts}.json`,
					"application/json; charset=utf-8",
					JSON.stringify(this.items, null, 2)
				)
				useNotificationsStore().success(
					"Данные выгружены в JSON",
					"История метрик"
				)
				return
			}

			const escapeCsv = (v: unknown) => {
				const s = String(v ?? "")
				const needsQuotes = /[",\n\r]/.test(s)
				const escaped = s.replace(/"/g, '""')
				return needsQuotes ? `"${escaped}"` : escaped
			}

			const header = [
				"timestamp",
				"cpuPercent",
				"memoryPercent",
				"memoryUsedMB",
				"memoryTotalMB",
			]
			const rows = this.items.map(item =>
				[
					escapeCsv(item.timestamp),
					escapeCsv(item.cpuPercent),
					escapeCsv(item.memoryPercent),
					escapeCsv(item.memoryUsedMB),
					escapeCsv(item.memoryTotalMB),
				].join(",")
			)
			const csv = [header.join(","), ...rows].join("\n")
			this.downloadFile(
				`metrics_history_${ts}.csv`,
				"text/csv; charset=utf-8",
				csv
			)
			useNotificationsStore().success(
				"Данные выгружены в CSV",
				"История метрик"
			)
		},
		/**
		 * Внутренний метод загрузки без учета "loading"-гарда.
		 */
		async fetchHistory(): Promise<void> {
			const from = this.toApiDateTime(this.fromDate)
			const to = this.toApiDateTime(this.toDate)

			const data = await this.$services.metricsHistory.getHistory({
				from: from || undefined,
				to: to || undefined,
				limit: this.limit,
			})
			this.items = Array.isArray(data) ? data : []
		},

		/**
		 * Преобразовать значение "datetime-local" в формат для api/БД: "YYYY-MM-DD HH:mm:ss".
		 *
		 */
		toApiDateTime(value: string): string {
			const v = String(value ?? "").trim()
			if (!v) return ""
			const s = v.length >= 19 ? v.slice(0, 19) : v
			return s.replace("T", " ")
		},

		/**
		 * Загрузить историю метрик по текущим фильтрам.
		 * "if (this.loading) return" Блокирую выполнение, если загрузка уже идет,
		 * чтобы избежать повторных запросов.
		 */
		async load(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null

			try {
				await this.fetchHistory()
			} catch (e: any) {
				this.error = e?.message ?? "Ошибка загрузки истории метрик"
			} finally {
				this.loading = false
			}
		},

		/**
		 * Очистить историю метрик (вся таблица) и перезагрузить список.
		 * "if (this.loading) return" Блокирую выполнение, если загрузка уже идет,
		 * чтобы избежать повторных запросов.
		 */
		async clearAll(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null

			try {
				const res = await this.$services.metricsHistory.clearHistory()
				useNotificationsStore().success(
					res.message ?? "Метрики успешно очищены",
					"История метрик"
				)
				await this.fetchHistory()
			} catch (e: any) {
				const msg = e?.message ?? "Не удалось очистить метрики"
				this.error = msg
				useNotificationsStore().error(msg, "История метрик")
			} finally {
				this.loading = false
			}
		},
	},
})
