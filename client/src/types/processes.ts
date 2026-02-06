export type ProcessInfo = {
	pid: number
	name: string
	exe: string
	cmdline: string
	username: string
	status: string
	createTime: number
	parentPid: number
	cpuPercent: number
	memoryPercent: number
	memoryRss: number
	ports: number[]
}

export type KillProcessResponse = {
	pid: number
	message: string
	timestamp: string
}

/**
 * Сообщение статуса от WS, когда мониторинг выключен/включен на сервере.
 */
export type MonitoringStatusWsPayload = {
	monitoringEnabled: boolean
	message?: string
}
