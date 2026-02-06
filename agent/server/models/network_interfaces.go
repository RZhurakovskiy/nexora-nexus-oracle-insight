package models

type NetworkInterfaceStat struct {
	Name         string `json:"name"`
	BytesSent    uint64 `json:"bytesSent"`
	BytesRecv    uint64 `json:"bytesRecv"`
	PacketsSent  uint64 `json:"packetsSent"`
	PacketsRecv  uint64 `json:"packetsRecv"`
	ErrIn        uint64 `json:"errIn"`
	ErrOut       uint64 `json:"errOut"`
	DropIn       uint64 `json:"dropIn"`
	DropOut      uint64 `json:"dropOut"`
}


