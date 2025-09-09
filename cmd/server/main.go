package main

import (
	"ConsoleChat/internal/infrastructure"
	"ConsoleChat/internal/service"
	"log"
	"time"
)

func main() {
	rs := service.NewRoomService(50, 30*time.Second) // последние 50 сообщений, таймаут 30 секунд
	cs := service.NewChatService(rs)
	server := infrastructure.NewTCPServer(":9000", rs, cs)

	log.Println("Запуск сервера ConsoleChat на порту 9000")
	server.Start()
}
