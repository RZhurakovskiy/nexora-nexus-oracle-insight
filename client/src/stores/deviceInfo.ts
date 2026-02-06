import type { DeviceInfo } from "@/services/host-info/HostInfoService"
import { defineStore } from "pinia"

type DeviceInfoState = {
	data: DeviceInfo | null
	loading: boolean
	error: string | null
	initialized: boolean
}

export const useDeviceInfoStore = defineStore("deviceInfo", {
	state: (): DeviceInfoState => ({
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
		 * Принудительно обновить данные с сервера.
		 */
		async refresh(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null

			try {
				this.data = await this.$services.hostInfo.getDeviceInfo()
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить информацию о процессоре"
			} finally {
				this.loading = false
			}
		},
	},
})
