import type { SystemInfoResponse } from "@/types/systemInfo"
import { defineStore } from "pinia"

type SystemInfoState = {
	data: SystemInfoResponse | null
	loading: boolean
	error: string | null
	initialized: boolean
}

export const useSystemInfoStore = defineStore("systemInfo", {
	state: (): SystemInfoState => ({
		data: null,
		loading: false,
		error: null,
		initialized: false,
	}),

	actions: {
		/**
		 * Инициализация страницы.
		 *
		 * не вызывать API лишний раз без явной перезагрузки.
		 */
		async init(): Promise<void> {
			if (this.loading || this.initialized) return
			await this.refresh()
			this.initialized = this.data != null && !this.error
		},

		/**
		 * Принудительное обновление данных.
		 */
		async refresh(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null
			try {
				this.data = await this.$services.systemInfo.getSystemInfo()
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить информацию о системе"
			} finally {
				this.loading = false
			}
		},
	},
})
