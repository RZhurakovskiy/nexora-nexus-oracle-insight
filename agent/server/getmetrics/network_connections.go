/* Функции для получения информации о сетевых соединениях */
package getmetrics

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

/* Определяет протокол соединения на основе типа и семейства адресов */
/* Возвращает строку с названием протокола */
func protoFromConn(conn net.ConnectionStat) string {
	base := "net"
	switch conn.Type {
	case 1:
		base = "tcp"
	case 2:
		base = "udp"
	}
	if conn.Family == 10 {
		return base + "6"
	}
	return base
}

/* Получает список всех сетевых соединений с информацией о процессах и сортирует их */
/* Возвращает массив структур с информацией о соединениях или ошибку */
func GetNetworkConnections(kind string) ([]models.ListeningPort, error) {
	if kind == "" {
		kind = "all"
	}

	conns, err := net.Connections(kind)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить сетевые соединения (%s): %w", kind, err)
	}

	result := make([]models.ListeningPort, 0, len(conns))
	for _, conn := range conns {
		if conn.Laddr.Port == 0 || conn.Pid <= 0 {
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
			Protocol:   protoFromConn(conn),
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
		if result[i].Port != result[j].Port {
			return result[i].Port < result[j].Port
		}
		return result[i].Process < result[j].Process
	})

	return result, nil
}
