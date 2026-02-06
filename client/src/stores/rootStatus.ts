import { defineStore } from "pinia"

type RootStatusState = {
	rootStatus: boolean
	error: string
}

export const useRootStatusStore = defineStore("rootStatus", {
	state: (): RootStatusState => ({
		rootStatus: false,
		error: "",
	}),

	actions: {
		/**
		 * Принудительно обновить данные с сервера.
		 */
		async init(): Promise<void> {
			try {
				const response = await this.$services.rootService.getRootService()
				this.rootStatus = response.rootStatus
			} catch (e: any) {
				this.error =
					e?.message ?? "Не удалось получить информацию о root правах"
			}
		},
	},
})
