export type SystemInfoHost = {
  hostname: string
  os: string
  platform: string
  platformFamily: string
  platformVersion: string
  kernelVersion: string
  kernelArch: string
  uptimeSec: number
  bootTime: number
}

export type SystemInfoCPU = {
  modelName: string
  vendorId: string
  physicalCores: number
  logicalCores: number
}

export type SystemInfoMemory = {
  totalBytes: number
  availableBytes: number
  usedBytes: number
  usedPercent: number
  swapTotalBytes: number
  swapUsedBytes: number
  swapUsedPercent: number
}

export type SystemInfoLoad = {
  load1: number
  load5: number
  load15: number
}

export type SystemInfoDisk = {
  device: string
  mountpoint: string
  fstype: string
  totalBytes: number
  usedBytes: number
  freeBytes: number
  usedPercent: number
}

export type SystemInfoNetInterface = {
  name: string
  addrs: string[]
}

export type SystemInfoResponse = {
  host: SystemInfoHost
  cpu: SystemInfoCPU
  memory: SystemInfoMemory
  load: SystemInfoLoad | null
  disks: SystemInfoDisk[]
  interfaces: SystemInfoNetInterface[]
}


