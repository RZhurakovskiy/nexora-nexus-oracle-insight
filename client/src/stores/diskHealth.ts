import type { DiskHealthResponse } from "@/types/diskHealth"
import { defineStore } from "pinia"

type DiskHealthState = {
	data: DiskHealthResponse | null
	loading: boolean
	error: string | null
}

export const useDiskHealthStore = defineStore("diskHealth", {
	state: (): DiskHealthState => ({
		data: null,
		loading: false,
		error: null,
	}),

	actions: {
		async init(): Promise<void> {
			if (this.loading) return
			if (this.data) return
			await this.refresh()
		},

		/**
		 * Принудительно обновить данные.
		 */
		async refresh(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null
			try {
				this.data = await this.$services.diskHealth.getDiskHealth()
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить SMART/health"
			} finally {
				this.loading = false
			}
		},
	},
})
