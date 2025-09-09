package service

import (
	"ConsoleChat/internal/domain"
	"fmt"
	"sync"
)

type ChatService struct {
	roomService *RoomService
	mu          sync.Mutex
}

func NewChatService(rs *RoomService) *ChatService {
	return &ChatService{roomService: rs}
}

func (cs *ChatService) SendMessage(roomID, sender, content string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	room, ok := cs.roomService.GetRoom(roomID)
	if !ok {
		fmt.Println("Попытка отправки в несуществующую комнату:", roomID)
		return
	}

	// Добавляем сообщение в историю
	room.Messages = append(room.Messages, domain.Message{Sender: sender, Content: content})
	if len(room.Messages) > cs.roomService.MessageLimit {
		room.Messages = room.Messages[len(room.Messages)-cs.roomService.MessageLimit:]
	}

	// Рассылка всем клиентам кроме отправителя
	for _, client := range room.Clients {
		if client.Name != sender {
			client.Conn.Write([]byte(fmt.Sprintf("%s: %s\n", sender, content)))
		}
	}
}

// Системные сообщения
func (cs *ChatService) SendSystemMessage(roomID, message string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	room, ok := cs.roomService.GetRoom(roomID)
	if !ok {
		return
	}

	for _, client := range room.Clients {
		client.Conn.Write([]byte(message))
	}
}
