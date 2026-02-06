export type StartRecordingParams = {
	cpuThreshold: number
	ramThreshold: number
	duration: number
}

export type StartRecordingResponse = {
	success: boolean
	sessionId: number
	message: string
}

export type RecordingSession = {
	id: number
	cpuThreshold: number
	ramThreshold: number
	duration: number
	startedAt: string
	endTime: string
}

export type RecordingStatusResponse = {
	active: boolean
	session?: RecordingSession
}
