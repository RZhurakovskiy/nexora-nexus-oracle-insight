/* Функции для получения топ процессов по количеству сетевых соединений */
package getmetrics

import (
	"path/filepath"
	"sort"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

/* Получает топ процессов по количеству сетевых соединений с разбивкой по статусам */
/* Сортирует по приоритету ESTABLISHED, общее количество, LISTEN */
/* Возвращает массив структур с информацией о процессах или ошибку */
func GetTopNetworkProcesses(limit int) ([]models.NetworkProcessStat, error) {
	if limit <= 0 {
		limit = 20
	}

	conns, _ := net.Connections("all")
	connCount := map[int32]int{}
	listenCount := map[int32]int{}
	establishedCount := map[int32]int{}
	otherCount := map[int32]int{}
	for _, c := range conns {
		if c.Pid <= 0 {
			continue
		}
		connCount[c.Pid]++
		if c.Status == "LISTEN" {
			listenCount[c.Pid]++
		} else if c.Status == "ESTABLISHED" {
			establishedCount[c.Pid]++
		} else {
			otherCount[c.Pid]++
		}
	}

	stats := make([]models.NetworkProcessStat, 0, len(connCount))
	for pid, connsN := range connCount {
		if connsN == 0 {
			continue
		}

		name := "неизвестно"
		username := ""

		p, err := process.NewProcess(pid)
		if err == nil {
			if exe, err := p.Exe(); err == nil && exe != "" {
				name = filepath.Base(exe)
			} else if n, err := p.Name(); err == nil && n != "" {
				name = n
			}
			if u, err := p.Username(); err == nil {
				username = u
			}
		}

		stats = append(stats, models.NetworkProcessStat{
			PID:         pid,
			Process:     name,
			Username:    username,
			Connections: connsN,
			Listening:   listenCount[pid],
			Established: establishedCount[pid],
			OtherStates: otherCount[pid],
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Established != stats[j].Established {
			return stats[i].Established > stats[j].Established
		}
		if stats[i].Connections != stats[j].Connections {
			return stats[i].Connections > stats[j].Connections
		}
		if stats[i].Listening != stats[j].Listening {
			return stats[i].Listening > stats[j].Listening
		}
		return stats[i].PID < stats[j].PID
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}
	return stats, nil
}
