import { defineStore } from "pinia"

type HostInfoState = {
	username: string
	hostname: string
	hostInfoLoading: boolean
	hostInfoInitialized: boolean

	processname: string
	cores: number | null
	machineInfoLoading: boolean
	machineInfoInitialized: boolean

	error: string | null
}

export const useHostInfoStore = defineStore("host-info", {
	state: (): HostInfoState => ({
		username: "",
		hostname: "",
		hostInfoLoading: false,
		hostInfoInitialized: false,

		processname: "",
		cores: null,
		machineInfoLoading: false,
		machineInfoInitialized: false,

		error: null,
	}),
	getters: {
		displayUsername: state => state.username || "Не определено",
		displayHostname: state => state.hostname || "Не определено",
		displayProcessname: state => state.processname || "Не определено",
		displayCores: state => state.cores || "Не определено",
	},
	actions: {
		/**
		 * Загружает имя пользователя и имя хоста.
		 *
		 * если ошибка hostInfoInitialized останется false что бы потом повторить попытку при след. вызове
		 */
		async initHostInfo(): Promise<void> {
			if (this.hostInfoLoading || this.hostInfoInitialized) return
			this.hostInfoLoading = true
			this.error = null
			try {
				const res = await this.$services.hostInfo.getHostInfo()
				this.username = res.username
				this.hostname = res.hostname
				this.hostInfoInitialized = true
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить информацию хоста"
			} finally {
				this.hostInfoLoading = false
			}
		},

		/**
		 * Загружает информацию о железе машины (процессор, количество ядер).
		 *
		 * если ошибка machineInfoLoading останется false что бы потом повторить попытку при след. вызове
		 */
		async initMachineInfo(): Promise<void> {
			if (this.machineInfoLoading || this.machineInfoInitialized) return
			this.machineInfoLoading = true
			this.error = null
			try {
				const res = await this.$services.hostInfo.getMachineInfo()
				this.processname = res.processname
				this.cores = res.cores
				this.machineInfoInitialized = true
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось получить информацию о машине"
			} finally {
				this.machineInfoLoading = false
			}
		},
	},
})
