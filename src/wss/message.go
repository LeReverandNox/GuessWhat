package wss

type Coords struct {
	X         int
	Y         int
	Color     string
	Thickness string
}

type Action struct {
	Action string
}

type Message struct {
	Data   string
	Coords *Coords
}

// type Message struct {
// 	// the json tag means this will serialize as a lowercased field
// 	Message string `json:"message"`
// 	Type    string `json:type`
// 	Content string `json:content`
// }
