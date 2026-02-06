/* Сервисы для управления запуском и мониторингом процессов через TCP */
package services

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	runningProcesses = make(map[int32]*exec.Cmd)
	procMutex        = sync.RWMutex{}
)

var AllowedCommands = map[string]bool{
	"node":     true,
	"npm":      true,
	"python":   true,
	"python3":  true,
	"go":       true,
	"vite":     true,
	"bun":      true,
	"deno":     true,
	"my-app":   true,
	"./my-app": true,
}

/* Структура для хранения результата запуска процесса */
type ProcessLaunchResult struct {
	PID   int32
	Msg   string
	Error error
}

/* Извлекает номер порта из аргументов команды */
/* Возвращает строку с номером порта или пустую строку если порт не найден */
func extractPortFromArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	last := args[len(args)-1]
	if port, err := strconv.Atoi(last); err == nil && port >= 1024 && port <= 65535 {
		return strconv.Itoa(port)
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "--port=") {
			portStr := strings.TrimPrefix(arg, "--port=")
			if port, err := strconv.Atoi(portStr); err == nil && port >= 1024 && port <= 65535 {
				return portStr
			}
		}
		if arg == "-p" || arg == "--port" {
		}
	}
	return ""
}

/* Проверяет существование скрипта если первый аргумент содержит точку */
/* Возвращает ошибку если файл не найден */
func isValidScript(cwd string, args []string) error {
	if len(args) == 0 {
		return nil
	}
	firstArg := args[0]
	if strings.Contains(firstArg, ".") {
		fullPath := firstArg
		if cwd != "" {
			fullPath = filepath.Join(cwd, firstArg)
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("файл не найден: %s", fullPath)
		}
	}
	return nil
}

/* Запускает новый процесс с указанной командой, аргументами и рабочей директорией */
/* Проверяет разрешенность команды и валидность путей, возвращает результат запуска или ошибку */
func StartProcess(command, argsStr, cwd string) (*ProcessLaunchResult, error) {
	if command == "" {
		return nil, fmt.Errorf("поле 'command' обязательно")
	}

	if !AllowedCommands[command] {
		return nil, fmt.Errorf("команда '%s' не разрешена", command)
	}

	var args []string
	if argsStr != "" {
		args = strings.Fields(argsStr)
	}

	if cwd != "" {
		if !filepath.IsAbs(cwd) {
			return nil, fmt.Errorf("cwd должен быть абсолютным путём")
		}
		if _, err := os.Stat(cwd); os.IsNotExist(err) {
			return nil, fmt.Errorf("директория cwd не существует: %s", cwd)
		}
	}

	if err := isValidScript(cwd, args); err != nil {
		return nil, err
	}

	cmd := exec.Command(command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("не удалось запустить: %w", err)
	}

	pid := int32(cmd.Process.Pid)

	procMutex.Lock()
	runningProcesses[pid] = cmd
	procMutex.Unlock()

	port := extractPortFromArgs(args)

	go monitorProcess(cmd, pid, command, argsStr, cwd, port, &stdout, &stderr)

	var msg string
	if port != "" {
		msg = fmt.Sprintf("Процесс запущен (PID=%d). Ожидается сервер на порту %s. Статус проверяется в фоне.", pid, port)
	} else {
		msg = fmt.Sprintf("Процесс запущен (PID=%d). Статус проверяется в фоне.", pid)
	}

	return &ProcessLaunchResult{
		PID: pid,
		Msg: msg,
	}, nil
}

/* Мониторит запущенный процесс, проверяет доступность порта если указан и логирует результат завершения */
func monitorProcess(cmd *exec.Cmd, pid int32, command, argsStr, cwd, port string, stdout, stderr *strings.Builder) {
	defer func() {
		procMutex.Lock()
		delete(runningProcesses, pid)
		procMutex.Unlock()
	}()

	time.Sleep(500 * time.Millisecond)

	procMutex.RLock()
	_, stillRunning := runningProcesses[pid]
	procMutex.RUnlock()

	if !stillRunning {
		cmd.Wait()
		return
	}

	if port != "" {
		addr := "127.0.0.1:" + port
		conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
		if err != nil {
			log.Printf("Процесс PID=%d запущен, но порт %s недоступен: %v", pid, port, err)
		} else {
			conn.Close()
			log.Printf("Процесс PID=%d работает на порту %s", pid, port)
		}
	}

	err := cmd.Wait()
	exitStatus := "успешно"
	if err != nil {
		exitStatus = fmt.Sprintf("с ошибкой: %v", err)
	}

	log.Printf("Процесс завершён PID=%d: %s %s (cwd: %s) → %s", pid, command, argsStr, cwd, exitStatus)
	if stderr.Len() > 0 {
		log.Printf("stderr PID=%d: %s", pid, strings.TrimSpace(stderr.String()))
	}
	if stdout.Len() > 0 {
		log.Printf("stdout PID=%d: %s", pid, strings.TrimSpace(stdout.String()))
	}
}
