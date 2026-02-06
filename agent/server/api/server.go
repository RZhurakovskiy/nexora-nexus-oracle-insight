/* Запуск и управление HTTP сервером с поддержкой graceful shutdown */
package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RZhurakovskiy/agent/server/middleware"
	"github.com/RZhurakovskiy/agent/server/services"
)

/*
Инициализирует базу данных, настраивает роуты и запускает HTTP сервер на указанном порту

	Ожидает сигналы SIGINT или SIGTERM
*/
func StartServer(port string) {

	sqlDB, err := InitDB("./monitor.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer sqlDB.Close()

	services.SetDB(sqlDB)

	mux := http.NewServeMux()

	SetupRoutes(mux)

	handler := middleware.CorsMiddleware(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.SetFlags(0)
		log.Println("==========================================")
		log.Printf("Proctl Server v1.0")
		log.Println("   Системный мониторинг и управление")
		log.Println("------------------------------------------")
		log.Printf("Сервер запущен на: http://localhost:%s", port)
		log.Println("------------------------------------------")
		log.Println("Стек:")
		log.Println("   - Go 1.25.3")
		log.Println("   - gorilla/websocket")
		log.Println("   - gopsutil v4 (системные метрики)")
		log.Println("   - agent ядро")
		log.Println("==========================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\nПолучен сигнал завершения. Останавливаем сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
