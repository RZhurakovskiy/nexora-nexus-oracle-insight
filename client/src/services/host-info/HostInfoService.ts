import type { ApiClient } from "@/services/http/ApiClient"

export type HostInfoResponse = {
	username: string
	hostname: string
}

export type MachineResponse = {
	processname: string
	cores: number
}

export type DeviceInfo = {
	processor_name: string
	vendor: string
	physical_cores: number
	logical_processors: number
	frequency_mhz: number
	cache_size_kb: number
	supported_flags: string[]
	architecture: string
	family: string
	model: string
}

export class HostInfoService {
	private readonly api: ApiClient

	constructor(api: ApiClient) {
		this.api = api
	}

	/**
	 * Получить имя пользователя и хоста системы.
	 */
	getHostInfo(): Promise<HostInfoResponse> {
		return this.api.getJson<HostInfoResponse>("/api/get-host-username")
	}

	/**
	 * Получить информацию о железе машины.
	 */
	getMachineInfo(): Promise<MachineResponse> {
		return this.api.getJson<MachineResponse>("/api/get-device-info")
	}

	/**
	 * Получить подробную информацию о процессоре.
	 */
	getDeviceInfo(): Promise<DeviceInfo> {
		return this.api.getJson<DeviceInfo>("/api/get-device-info")
	}
}
