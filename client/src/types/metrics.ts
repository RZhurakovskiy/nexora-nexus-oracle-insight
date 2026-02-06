export type CpuWsPayload = {
  cpu?: number
  timestamp?: string
  monitoringEnabled?: boolean
  message?: string
  error?: string
}

export type MemoryWsPayload = {
  memoryUsage?: number
  usedMB?: number
  totalMemory?: number
  timestamp?: string
  monitoringEnabled?: boolean
  message?: string
  error?: string
}

export type CpuSample = {
  cpu: number
  timestamp: string
}

export type MemorySample = {
  memoryUsage: number
  usedMB: number
  totalMemoryMB: number
  timestamp: string
}


