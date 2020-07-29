package messages

import (
	"errors"
	"time"

	"github.com/will-slack/apiserver/models/users"
	"gopkg.in/mgo.v2/bson"
)

//MessageID defines the type for message IDs
type MessageID interface{}

//Message represents a message in the database
type Message struct {
	ID        MessageID    `json:"id" bson:"_id"`
	ChannelID ChannelID    `json:"channelID"`
	Body      string       `json:"body"`
	CreatedAt time.Time    `json:"createdAt"`
	CreatorID users.UserID `json:"creatorID"`
	EditedAt  time.Time    `json:"editedAt"`
	FirstName string       `json:"firstname"`
	LastName  string       `json:"lastname"`
}

//NewMessage represents a new message created in a chat
type NewMessage struct {
	ChannelID string `json:"channelID"`
	Body      string `json:"body"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//MessageUpdates represents updates a channel creator can make
type MessageUpdates struct {
	Body string
}

//Validate validates a new message
func (nm *NewMessage) Validate() error {
	if len(nm.ChannelID) == 0 || len(nm.Body) == 0 {
		return errors.New("Values not appropriate")
	}
	return nil
}

//ToMessage converts the NewMessage to a Message
func (nm *NewMessage) ToMessage(creator *users.User) (*Message, error) {
	if id, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(id)
	}
	message := &Message{
		ChannelID: nm.ChannelID,
		Body:      nm.Body,
		CreatedAt: time.Now(),
		CreatorID: creator.ID,
		EditedAt:  time.Now(),
		FirstName: nm.FirstName,
		LastName:  nm.LastName,
	}
	return message, nil
}
