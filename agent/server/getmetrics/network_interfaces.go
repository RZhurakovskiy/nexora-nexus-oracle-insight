package getmetrics

import (
	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
)

func GetNetworkInterfacesIO() ([]models.NetworkInterfaceStat, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	result := make([]models.NetworkInterfaceStat, 0, len(counters))
	for _, c := range counters {
		result = append(result, models.NetworkInterfaceStat{
			Name:        c.Name,
			BytesSent:   c.BytesSent,
			BytesRecv:   c.BytesRecv,
			PacketsSent: c.PacketsSent,
			PacketsRecv: c.PacketsRecv,
			ErrIn:       c.Errin,
			ErrOut:      c.Errout,
			DropIn:      c.Dropin,
			DropOut:     c.Dropout,
		})
	}
	return result, nil
}
