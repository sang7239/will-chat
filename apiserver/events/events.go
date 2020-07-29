package events

// Event defines a event that is transmitted via websocket to a client
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
