export type NetworkConnection = {
	port: number
	protocol: string
	pid: number
	process: string
	status: string
	localAddr: string
	remoteAddr: string
}

export type NetworkProcessStat = {
	pid: number
	process: string
	username: string
	connections: number
	listening: number
	established: number
	otherStates: number
}

export type NetworkInterfaceStat = {
	name: string
	bytesSent: number
	bytesRecv: number
	packetsSent: number
	packetsRecv: number
	errIn: number
	errOut: number
	dropIn: number
	dropOut: number
}
