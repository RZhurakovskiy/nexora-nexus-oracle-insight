export type DiskHealthDevice = {
	device: string
	model: string
	serial: string
	busType: string
	smartAvailable: boolean
	smartEnabled: boolean
	smartPassed: boolean | null
	temperatureC: number | null
	powerOnHours: number | null
	unsafeShutdowns: number | null
	nvmePercentUsed: number | null
	warnings: string[]
}

export type DiskHealthResponse = {
	supported: boolean
	message: string
	devices: DiskHealthDevice[]
}
