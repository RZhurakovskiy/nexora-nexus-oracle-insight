export type ListeningPort = {
	port: number
	protocol: string
	pid: number
	process: string
	status: string
	localAddr: string
	remoteAddr: string
}

export type StartProcessRequest = {
	command: string
	args: string
	cwd?: string
	timestamp: string
}

export type StartProcessResponse = {
	pid: number
	command: string
	args: string
	cwd?: string
	msg: string
}
