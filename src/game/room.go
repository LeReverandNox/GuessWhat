package game

type Room struct {
	Name    string
	Clients []*Client
	Drawer  string
	Image   string
}

func NewRoom(name string) *Room {
	room := Room{}
	room.Name = name
	room.Clients = make([]*Client, 0)

	return &room
}
