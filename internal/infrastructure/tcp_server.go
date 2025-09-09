package infrastructure

import (
	"ConsoleChat/internal/domain"
	"ConsoleChat/internal/service"
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type TCPServer struct {
	Addr        string
	RoomService *service.RoomService
	ChatService *service.ChatService
}

func NewTCPServer(addr string, rs *service.RoomService, cs *service.ChatService) *TCPServer {
	return &TCPServer{Addr: addr, RoomService: rs, ChatService: cs}
}

func (s *TCPServer) Start() {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	log.Println("Server listening on", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("Введите ваше имя:\n"))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	conn.Write([]byte("Введите 'create' для новой комнаты или UUID существующей:\n"))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var roomID string
	if choice == "create" {
		roomID = generateUUID()
		s.RoomService.CreateRoom(roomID)
		conn.Write([]byte(fmt.Sprintf("Комната создана! Ваш UUID: %s\n", roomID)))
	} else {
		roomID = choice
		if _, ok := s.RoomService.GetRoom(roomID); !ok {
			s.RoomService.CreateRoom(roomID)
			conn.Write([]byte(fmt.Sprintf("Создана новая комната: %s\n", roomID)))
		} else {
			conn.Write([]byte(fmt.Sprintf("Подключились к комнате %s\n", roomID)))
		}
	}

	client := &domain.Client{Conn: conn, Name: name}
	s.RoomService.AddClient(roomID, client)

	// Отправка последних сообщений
	if room, ok := s.RoomService.GetRoom(roomID); ok {
		for _, msg := range room.Messages {
			conn.Write([]byte(fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content)))
		}
	}

	// Сообщение о подключении
	s.ChatService.SendSystemMessage(roomID, fmt.Sprintf("[Система] %s подключился\n", name))

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			s.RoomService.RemoveClient(roomID, name)
			s.ChatService.SendSystemMessage(roomID, fmt.Sprintf("[Система] %s отключился\n", name))
			return
		}
		msg = strings.TrimSpace(msg)
		s.ChatService.SendMessage(roomID, name, msg)
	}
}

func generateUUID() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
