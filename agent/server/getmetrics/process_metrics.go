/* Функции для получения информации о процессах системы */
package getmetrics

import (
	"sync"

	"github.com/RZhurakovskiy/agent/server/models"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

/* Получает информацию о всех процессах системы с использованием параллельной обработки */
/* Собирает метрики cpu, памяти, сетевые порты и другую информацию о каждом процессе */
/* Возвращает массив структур с информацией о процессах или ошибку */
func UsageProcess(allConnections []net.ConnectionStat) ([]models.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	connectionsByPID := make(map[int32][]net.ConnectionStat)
	for _, conn := range allConnections {
		if conn.Pid > 0 {
			connectionsByPID[conn.Pid] = append(connectionsByPID[conn.Pid], conn)
		}
	}

	result := make([]models.ProcessInfo, 0, len(procs))
	resultMutex := sync.Mutex{}

	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, p := range procs {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(proc *process.Process) {
			defer wg.Done()
			defer func() { <-semaphore }()

			info := models.ProcessInfo{
				PID: proc.Pid,
			}

			if name, err := proc.Name(); err == nil {
				info.Name = name
			} else {
				info.Name = "неизвестно"
			}

			if exe, err := proc.Exe(); err == nil {
				info.Exe = exe
			}

			if cmdline, err := proc.Cmdline(); err == nil {
				info.Cmdline = cmdline
			}

			if username, err := proc.Username(); err == nil {
				info.Username = username
			}

			if status, err := proc.Status(); err == nil && len(status) > 0 {
				info.Status = status[0]
			}

			if createTime, err := proc.CreateTime(); err == nil {
				info.CreateTime = createTime
			}

			if parentProc, err := proc.Parent(); err == nil && parentProc != nil {
				info.ParentPID = parentProc.Pid
			}

			if cpu, err := proc.CPUPercent(); err == nil {
				info.CPUPercent = cpu
			}

			if mem, err := proc.MemoryPercent(); err == nil {
				info.MemoryPercent = float64(mem)
			}

			if memInfo, err := proc.MemoryInfo(); err == nil {
				info.MemoryRSS = memInfo.RSS
			}

			info.Ports = make([]uint32, 0)
			if connections, ok := connectionsByPID[proc.Pid]; ok {
				portMap := make(map[uint32]bool)
				for _, conn := range connections {
					if conn.Laddr.Port > 0 {
						portMap[conn.Laddr.Port] = true
					}
				}
				info.Ports = make([]uint32, 0, len(portMap))
				for port := range portMap {
					info.Ports = append(info.Ports, port)
				}
			}

			resultMutex.Lock()
			result = append(result, info)
			resultMutex.Unlock()
		}(p)
	}

	wg.Wait()

	return result, nil
}

/* Получает информацию об ограниченном количестве процессов без параллельной обработки */
/* Используется когда нужно быстро получить информацию о небольшом количестве процессов */
/* Возвращает массив структур с информацией о процессах или ошибку */
func UsageProcessLimited(limit int, allConnections []net.ConnectionStat) ([]models.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(procs) > limit {
		procs = procs[:limit]
	}

	connectionsByPID := make(map[int32][]net.ConnectionStat)
	for _, conn := range allConnections {
		if conn.Pid > 0 {
			connectionsByPID[conn.Pid] = append(connectionsByPID[conn.Pid], conn)
		}
	}

	result := make([]models.ProcessInfo, 0, len(procs))

	for _, p := range procs {
		info := models.ProcessInfo{PID: p.Pid}

		if name, err := p.Name(); err == nil {
			info.Name = name
		} else {
			info.Name = "неизвестно"
		}

		if exe, err := p.Exe(); err == nil {
			info.Exe = exe
		}
		if cmdline, err := p.Cmdline(); err == nil {
			info.Cmdline = cmdline
		}
		if username, err := p.Username(); err == nil {
			info.Username = username
		}
		if status, err := p.Status(); err == nil && len(status) > 0 {
			info.Status = status[0]
		}
		if createTime, err := p.CreateTime(); err == nil {
			info.CreateTime = createTime
		}
		if parentProc, err := p.Parent(); err == nil && parentProc != nil {
			info.ParentPID = parentProc.Pid
		}
		if cpu, err := p.CPUPercent(); err == nil {
			info.CPUPercent = cpu
		}
		if mem, err := p.MemoryPercent(); err == nil {
			info.MemoryPercent = float64(mem)
		}
		if memInfo, err := p.MemoryInfo(); err == nil {
			info.MemoryRSS = memInfo.RSS
		}

		info.Ports = make([]uint32, 0)
		if connections, ok := connectionsByPID[p.Pid]; ok {
			portMap := make(map[uint32]bool)
			for _, conn := range connections {
				if conn.Laddr.Port > 0 {
					portMap[conn.Laddr.Port] = true
				}
			}
			info.Ports = make([]uint32, 0, len(portMap))
			for port := range portMap {
				info.Ports = append(info.Ports, port)
			}
		}

		result = append(result, info)
	}

	return result, nil
}
