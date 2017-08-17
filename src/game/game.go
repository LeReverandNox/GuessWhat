package game

type Game struct {
	Clients  []*Client
	Rooms    []*Room
	Messages []*Message
}

func NewGame() *Game {
	game := Game{}
	game.Clients = make([]*Client, 0)
	game.Rooms = make([]*Room, 0)
	game.Messages = make([]*Message, 0)

	return &game
}
