package messages

import "github.com/will-slack/apiserver/models/users"

//Store represents an abstract store for model.Message objects.
//This interface is used by the HTTP handlers to provide various functions
//such as getting all channels, inserting, updating, deleting channels etc.
//This interface can be implemented for any persistence database chosen (e.g MongoDB, PosgreSQL)
type Store interface {

	//Get all channels a given user is allowed to see
	//including channels the user is a member of, as well as all public channels
	GetAll(*users.User) ([]*Channel, error)

	//Insert inserts a new channel
	InsertChannel(newChannel *NewChannel, creator *users.User) (*Channel, error)

	//GetMostRecentMessages gets the most recent N messages posted to a particular channel
	GetMostRecentMessages(user *users.User, channelID interface{}, n int) ([]*Message, error)

	//GetChannelByID retrieves a channel with the given channelID
	GetChannelByID(channelID interface{}) (*Channel, error)

	//Update updates a channel's Name and Description
	UpdateChannel(updates *ChannelUpdates, channelID interface{}, user *users.User) error

	//Delete deletes a channel, as well as all messages posted to that channel
	DeleteChannel(channelID interface{}, user *users.User) error

	//AddUser adds a user to a channel's Members list
	AddUserToChannel(userID interface{}, channelID interface{}, creatorID interface{}) error

	//RemoveUserFromChannel removes a suer from a channel's Members list
	RemoveUserFromChannel(userID interface{}, channelID interface{}, creatorID interface{}) error

	//GetMessageByID retrieves a message with the given messageID
	GetMessageByID(MessageID interface{}) (*Message, error)

	//InsertMessage inserts a new message
	InsertMessage(message *NewMessage, user *users.User) (*Message, error)

	//UpdateMessage updates an existing message
	UpdateMessage(updates *MessageUpdates, messageID interface{}, creator *users.User) error

	//DeleteMessage deletes a message
	DeleteMessage(messageID interface{}, user *users.User) error
}
