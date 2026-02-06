import { CLIENT_VERSION } from "@/config/version"
import { defineStore } from "pinia"

/**
 * Версии клиента и сервера для отображения в ui.
 */
export const useVersionStore = defineStore("version", {
	state: () => ({
		/**
		 * Версия клиента
		 */
		clientVersion: CLIENT_VERSION,
		/**
		 * Версия сервера
		 */
		serverVersion: null as string | null,
		/**
		 * Загрузка версии
		 */
		loading: false,
		/**
		 * Ошибка загрузки
		 */
		error: null as string | null,
		/**
		 * Чтобы вызывать API один раз при переключении страниц.
		 */
		initialized: false,
	}),
	actions: {
		async init(): Promise<void> {
			if (this.initialized) return
			this.initialized = true
			await this.refresh()
		},

		/**
		 * Повторно запрашивает версию сервера.
		 */
		async refresh(): Promise<void> {
			this.loading = true
			this.error = null
			try {
				const res = await this.$services.version.getVersion()
				this.serverVersion = res.serverVersion
			} catch (e: any) {
				this.serverVersion = null
				this.error =
					typeof e?.message === "string"
						? e.message
						: "Не удалось получить версию сервера"
			} finally {
				this.loading = false
			}
		},
	},
})
