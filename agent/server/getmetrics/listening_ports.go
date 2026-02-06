package getmetrics

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

func GetListeningPorts() ([]models.ListeningPort, error) {
	conns, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить TCP-соединения: %w", err)
	}

	var result []models.ListeningPort

	for _, conn := range conns {
		if conn.Laddr.Port == 0 || conn.Pid <= 0 {
			continue
		}

		proto := ""
		switch conn.Family {
		case 2:
			proto = "tcp"
		case 10:
			proto = "tcp6"
		default:
			continue
		}

		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)

		remoteAddr := "0.0.0.0:0"
		if conn.Raddr.IP != "" && conn.Raddr.Port != 0 {
			remoteAddr = fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)
		} else if conn.Status == "LISTEN" {
			remoteAddr = "-"
		}

		processName := "неизвестно"
		p, err := process.NewProcess(conn.Pid)
		if err == nil {
			if exe, err := p.Exe(); err == nil && exe != "" {
				processName = filepath.Base(exe)
			} else if name, err := p.Name(); err == nil && name != "" {
				processName = name
			}
		}

		result = append(result, models.ListeningPort{
			Port:       conn.Laddr.Port,
			Protocol:   proto,
			PID:        conn.Pid,
			Process:    processName,
			Status:     conn.Status,
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Status == "LISTEN" && result[j].Status != "LISTEN" {
			return true
		}
		if result[i].Status != "LISTEN" && result[j].Status == "LISTEN" {
			return false
		}
		return result[i].Port < result[j].Port
	})

	return result, nil
}
