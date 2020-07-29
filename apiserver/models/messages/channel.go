package messages

import (
	"errors"
	"time"

	"github.com/will-slack/apiserver/models/users"
	"gopkg.in/mgo.v2/bson"
)

//ChannelID defines the type for channel IDs
type ChannelID interface{}

//Channel represents a chat channel in the database
type Channel struct {
	ID          ChannelID      `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	CreatorID   users.UserID   `json:"creatorID"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

//NewChannel represents a new channel created
type NewChannel struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Members     []users.UserID `json:"members,omitempty"`
	Private     bool           `json:"private,omitempty"`
}

//ChannelUpdates represents updates a user can make to a channel
type ChannelUpdates struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//Validate validates the new channel
func (nc *NewChannel) Validate() error {
	if len(nc.Name) == 0 {
		return errors.New("No channel name given")
	}
	return nil
}

//ToChannel converts the Newchannel to a channel
func (nc *NewChannel) ToChannel(creator *users.User) (*Channel, error) {
	if id, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(id)
	}
	channel := &Channel{
		Name:        nc.Name,
		Description: nc.Description,
		CreatedAt:   time.Now(),
		CreatorID:   creator.ID,
		Members:     nc.Members,
		Private:     nc.Private,
	}
	var members []users.UserID
	if nc.Members == nil {
		members = make([]users.UserID, 0, 10)
		members = append(members, creator.ID)
	} else {
		members = nc.Members
	}
	channel.Members = members
	return channel, nil
}
