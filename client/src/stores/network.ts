import type { ConnectionsKind } from "@/services/network/NetworkService"
import type {
	NetworkConnection,
	NetworkInterfaceStat,
	NetworkProcessStat,
} from "@/types/network"
import { defineStore } from "pinia"

type NetworkState = {
	connections: NetworkConnection[]
	topProcesses: NetworkProcessStat[]
	interfaces: NetworkInterfaceStat[]

	loading: boolean
	error: string | null

	kind: ConnectionsKind
	showOnlyListen: boolean

	filter: {
		query: string
		status: string
	}

	autoRefresh: boolean
}

let autoTimer: number | null = null

export const useNetworkStore = defineStore("network", {
	state: (): NetworkState => ({
		connections: [],
		topProcesses: [],
		interfaces: [],
		loading: false,
		error: null,
		kind: "all",
		showOnlyListen: false,
		filter: { query: "", status: "" },
		autoRefresh: false,
	}),

	getters: {
		filteredConnections(state): NetworkConnection[] {
			const q = state.filter.query.trim().toLowerCase()
			const status = state.filter.status.trim().toUpperCase()
			let list = state.connections
			if (state.showOnlyListen) {
				list = list.filter(
					x => String(x.status || "").toUpperCase() === "LISTEN"
				)
			}
			if (status) {
				list = list.filter(connection => String(connection.status || "").toUpperCase() === status)
			}
			if (q) {
				list = list.filter(connection => {
					const hay =
						`${connection.pid} ${connection.process} ${connection.protocol} ${connection.localAddr} ${connection.remoteAddr} ${connection.status} ${connection.port}`.toLowerCase()
					return hay.includes(q)
				})
			}
			return list
		},
	},

	actions: {
		/**
		 * Один раз загрузить данные.
		 * "if (this.loading) return" Блокирую выполнение, если загрузка уже идет,
		 * чтобы избежать повторных запросов.
		 */
		async init(): Promise<void> {
			if (this.loading) return
			if (this.connections.length > 0) return
			await this.refreshAll()
		},

		/**
		 * Обновить все соединения, топ процессов, интерфейсы.
		 */
		async refreshAll(): Promise<void> {
			if (this.loading) return
			this.loading = true
			this.error = null
			try {
				const [conns, top, ifaces] = await Promise.all([
					this.$services.network.getConnections(this.kind),
					this.$services.network.getTopProcesses(20),
					this.$services.network.getInterfaces(),
				])
				this.connections = Array.isArray(conns) ? conns : []
				this.topProcesses = Array.isArray(top) ? top : []
				this.interfaces = Array.isArray(ifaces) ? ifaces : []
			} catch (e: any) {
				this.error = e?.message ?? "Не удалось загрузить сетевую статистику"
			} finally {
				this.loading = false
			}
		},

		/**
		 * Включить/выключить автообновление.
		 */
		setAutoRefresh(enabled: boolean): void {
			this.autoRefresh = enabled
			if (autoTimer != null) window.clearInterval(autoTimer)
			autoTimer = null
			if (!enabled) return
			autoTimer = window.setInterval(() => {
				this.refreshAll()
			}, 5000)
		},
	},
})
