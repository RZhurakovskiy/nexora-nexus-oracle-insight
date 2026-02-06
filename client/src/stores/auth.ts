import { useNotificationsStore } from "@/stores/notifications"
import { defineStore } from "pinia"

type AuthState = {
	isAuthenticated: boolean
	loading: boolean
}

export const useAuthStore = defineStore("auth", {
	state: (): AuthState => ({
		isAuthenticated: false,
		loading: false,
	}),

	actions: {
		/**
		 * Проверить статус авторизации.
		 */
		checkAuth(): void {
			this.isAuthenticated = this.$services.auth.isAuthenticated()
		},

		/**
		 * Выполнить вход в систему.
		 */
		async login(login: string, password: string): Promise<void> {
			this.loading = true
			const notifications = useNotificationsStore()

			try {
				await this.$services.auth.login(login, password)
				this.isAuthenticated = true
				notifications.success("Успешный вход в систему")
			} catch (error) {
				const message =
					error instanceof Error ? error.message : "Ошибка авторизации"
				notifications.error(message, "Ошибка входа")
				throw error
			} finally {
				this.loading = false
			}
		},

		/**
		 * Выйти из системы.
		 */
		logout(): void {
			this.$services.auth.logout()
			this.isAuthenticated = false
			const notifications = useNotificationsStore()
			notifications.info("Вы вышли из системы")
		},

		/**
		 * Полный выход из системы с остановкой всех процессов.
		 */
		async logoutWithCleanup(): Promise<void> {
			const metrics = (await import("@/stores/metrics")).useMetricsStore()
			const processes = (await import("@/stores/processes")).useProcessesStore()
			const recording = (await import("@/stores/recording")).useRecordingStore()

			if (metrics.monitoringEnabled) {
				await metrics.stopMonitoring()
			}
			metrics.disconnectSockets()

			if (processes.socketsConnected) {
				processes.disconnectSockets()
			}

			if (recording.active) {
				await recording.stop()
			}
			recording.stopTick()

			this.logout()
		},
	},
})

