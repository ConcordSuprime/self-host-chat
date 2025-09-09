package service

import (
	"ConsoleChat/internal/domain"
	"fmt"
	"sync"
	"time"
)

type RoomService struct {
	rooms        map[string]*domain.Room
	mu           sync.Mutex
	MessageLimit int
	Timeout      time.Duration
}

func NewRoomService(messageLimit int, timeout time.Duration) *RoomService {
	return &RoomService{
		rooms:        make(map[string]*domain.Room),
		MessageLimit: messageLimit,
		Timeout:      timeout,
	}
}

func (rs *RoomService) CreateRoom(id string) *domain.Room {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	room := &domain.Room{
		ID:      id,
		Clients: make(map[string]*domain.Client),
	}
	rs.rooms[id] = room
	go rs.monitorRoom(room)
	return room
}

func (rs *RoomService) GetRoom(id string) (*domain.Room, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	room, ok := rs.rooms[id]
	return room, ok
}

func (rs *RoomService) AddClient(roomID string, client *domain.Client) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if room, ok := rs.rooms[roomID]; ok {
		room.Clients[client.Name] = client
	}
}

func (rs *RoomService) RemoveClient(roomID, clientName string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if room, ok := rs.rooms[roomID]; ok {
		delete(room.Clients, clientName)
	}
}

func (rs *RoomService) monitorRoom(room *domain.Room) {
	for {
		time.Sleep(rs.Timeout)
		rs.mu.Lock()
		if len(room.Clients) == 0 {
			delete(rs.rooms, room.ID)
			rs.mu.Unlock()
			fmt.Println("Комната", room.ID, "удалена из-за простоя")
			return
		}
		rs.mu.Unlock()
	}
}
