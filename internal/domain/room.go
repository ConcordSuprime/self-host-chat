package domain

type Room struct {
	ID       string
	Clients  map[string]*Client
	Messages []Message
}
