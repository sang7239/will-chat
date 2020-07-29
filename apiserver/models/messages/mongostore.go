package messages

import (
	"errors"

	"github.com/will-slack/apiserver/models/users"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const defaultAddr = "127.0.0.0.1:27017"

//MongoStore is an implementation of MessageStore
//backed by MongoDB
type MongoStore struct {
	Session           *mgo.Session
	Database          string
	MessageCollection string
	ChannelCollection string
}

//NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, databaseName string) *MongoStore {
	var err error
	if session == nil {
		session, err = mgo.Dial(defaultAddr)
	}
	if err != nil {
		return nil
	}
	if databaseName == "" {
		databaseName = "production"
	}
	return &MongoStore{
		Session:           session,
		Database:          databaseName,
		MessageCollection: "messages",
		ChannelCollection: "channels",
	}
}

//GetAll channels a given user is allowed to see
//including channels the user is a member of, as well as all public channels
func (ms *MongoStore) GetAll(user *users.User) ([]*Channel, error) {
	channels := []*Channel{}
	coll := ms.Session.DB(ms.Database).C(ms.ChannelCollection)
	if err := coll.Find(bson.M{"$or": []bson.M{{"members": user.ID}, {"private": false}}}).All(&channels); err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return channels, nil
}

//InsertChannel inserts a new channel
func (ms *MongoStore) InsertChannel(newChannel *NewChannel, creator *users.User) (*Channel, error) {
	if id, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(id)
	}
	channel, err := newChannel.ToChannel(creator)
	if err != nil {
		return nil, err
	}
	channel.ID = bson.NewObjectId()
	coll := ms.Session.DB(ms.Database).C(ms.ChannelCollection)
	if err := coll.Insert(channel); err != nil {
		return nil, err
	}
	return channel, nil
}

//GetMostRecentMessages gets the most recent N messages posted to a particular channel
func (ms *MongoStore) GetMostRecentMessages(user *users.User, channelID interface{}, n int) ([]*Message, error) {
	if id, ok := user.ID.(string); ok {
		user.ID = bson.ObjectIdHex(id)
	}
	coll := ms.Session.DB(ms.Database).C(ms.ChannelCollection)
	if err := coll.Find(bson.M{"$or": []bson.M{{"members": user.ID}, {"private": false}}}).One(nil); err != nil {
		return nil, errors.New("Not Authorized to view messages in this channel")
	}
	messages := []*Message{}
	coll = ms.Session.DB(ms.Database).C(ms.MessageCollection)
	if err := coll.Find(bson.M{"channelid": channelID}).Sort("createdat").Limit(n).All(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

//GetChannelByID retrieves a channel with the given channelID
func (ms *MongoStore) GetChannelByID(channelID interface{}) (*Channel, error) {
	if id, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(id)
	}
	coll := ms.Session.DB(ms.Database).C(ms.ChannelCollection)
	channel := &Channel{}
	if err := coll.Find(bson.M{"_id": channelID}).One(&channel); err != nil {
		return nil, err
	}
	return channel, nil
}

//UpdateChannel updates a channel's Name and Description
func (ms *MongoStore) UpdateChannel(updates *ChannelUpdates, channelID interface{}, user *users.User) error {
	if id, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(id)
	}
	if !ms.isChannelCreator(user, channelID) {
		return errors.New("Not Authorized to Update Channel")
	}
	coll := ms.Session.DB(ms.Database).C(ms.ChannelCollection)
	if err := coll.UpdateId(channelID, bson.M{"$set": updates}); err != nil {
		return err
	}
	return nil
}

//DeleteChannel deletes a channel, as well as all messages posted to that channel
func (ms *MongoStore) DeleteChannel(channelID interface{}, user *users.User) error {
	coll := ms.Session.DB(ms.Database).C(ms.MessageCollection)
	if _, err := coll.RemoveAll(bson.M{"channelid": channelID}); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}
	if id, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(id)
	}
	if !ms.isChannelCreator(user, channelID) {
		return errors.New("Not Authorized to Delete Channel")
	}
	if err := ms.Session.DB(ms.Database).C(ms.ChannelCollection).Remove(bson.M{"_id": channelID}); err != nil {
		return err
	}
	return nil
}

//AddUserToChannel adds a user to a channel's Members list
func (ms *MongoStore) AddUserToChannel(userID interface{}, channelID interface{}, creatorID interface{}) error {
	if id, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(id)
	}
	if id, ok := creatorID.(string); ok {
		creatorID = bson.ObjectIdHex(id)
	}
	if _, err := ms.Session.DB(ms.Database).C(ms.ChannelCollection).UpsertId(channelID, bson.M{"$addToSet": bson.M{"members": userID}}); err != nil {
		return err
	}
	return nil
}

//RemoveUserFromChannel removes a user from a channel's Members list
func (ms *MongoStore) RemoveUserFromChannel(userID interface{}, channelID interface{}, creatorID interface{}) error {
	if id, ok := channelID.(string); ok {
		channelID = bson.ObjectIdHex(id)
	}
	if id, ok := creatorID.(string); ok {
		creatorID = bson.ObjectIdHex(id)
	}
	if err := ms.Session.DB(ms.Database).C(ms.ChannelCollection).UpdateId(channelID, bson.M{"$pull": bson.M{"members": userID}}); err != nil {
		return err
	}
	return nil
}

//GetMessageByID retrieves a message with the given messageID
func (ms *MongoStore) GetMessageByID(messageID interface{}) (*Message, error) {
	if id, ok := messageID.(string); ok {
		messageID = bson.ObjectIdHex(id)
	}
	coll := ms.Session.DB(ms.Database).C(ms.MessageCollection)
	message := &Message{}
	if err := coll.Find(bson.M{"_id": messageID}).One(&message); err != nil {
		return nil, err
	}
	return message, nil
}

//InsertMessage inserts a new message
func (ms *MongoStore) InsertMessage(newMessage *NewMessage, creator *users.User) (*Message, error) {
	if id, ok := creator.ID.(string); ok {
		creator.ID = bson.ObjectIdHex(id)
	}
	message, err := newMessage.ToMessage(creator)
	if err != nil {
		return nil, err
	}
	message.ID = bson.NewObjectId()
	coll := ms.Session.DB(ms.Database).C(ms.MessageCollection)
	if err := coll.Insert(message); err != nil {
		return nil, err
	}
	return message, nil
}

//UpdateMessage updates an existing message
func (ms *MongoStore) UpdateMessage(updates *MessageUpdates, messageID interface{}, creator *users.User) error {
	if id, ok := messageID.(string); ok {
		messageID = bson.ObjectIdHex(id)
	}
	if !ms.isMessageCreator(creator, messageID) {
		return errors.New("Not Authorized to Update Message")
	}
	coll := ms.Session.DB(ms.Database).C(ms.MessageCollection)
	if err := coll.UpdateId(messageID, bson.M{"$set": updates}); err != nil {
		return err
	}
	return nil
}

//DeleteMessage deletes a message
func (ms *MongoStore) DeleteMessage(messageID interface{}, user *users.User) error {
	if id, ok := messageID.(string); ok {
		messageID = bson.ObjectIdHex(id)
	}
	if !ms.isMessageCreator(user, messageID) {
		return errors.New("Not Authorized to Delete Message")
	}
	coll := ms.Session.DB(ms.Database).C(ms.MessageCollection)
	if err := coll.Remove(bson.M{"_id": messageID}); err != nil {
		return err
	}
	return nil
}

func (ms *MongoStore) isChannelCreator(user *users.User, channelID interface{}) bool {
	channel := &Channel{}
	if err := ms.Session.DB(ms.Database).C(ms.ChannelCollection).Find(bson.M{"_id": channelID}).One(channel); err != nil {
		return false
	}
	if cid, ok := channel.CreatorID.(bson.ObjectId); ok {
		return cid.Hex() == user.ID
	}
	return false
}

func (ms *MongoStore) isMessageCreator(user *users.User, messageID interface{}) bool {
	message := &Message{}
	if err := ms.Session.DB(ms.Database).C(ms.MessageCollection).Find(bson.M{"_id": messageID}).One(message); err != nil {
		return false
	}
	if mid, ok := message.CreatorID.(bson.ObjectId); ok {
		return mid.Hex() == user.ID
	}
	return false
}
