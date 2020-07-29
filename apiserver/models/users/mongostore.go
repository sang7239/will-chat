package users

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const defaultAddr = "127.0.0.1:27017"

//MongoStore is an implementation of UserStore
//backed by MongoDB
type MongoStore struct {
	Session    *mgo.Session
	Database   string
	Collection string
}

//NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, databaseName string, collectionName string) *MongoStore {
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
	if collectionName == "" {
		collectionName = "users"
	}
	return &MongoStore{
		Session:    session,
		Database:   databaseName,
		Collection: collectionName,
	}
}

//GetAll returns all users
func (ms *MongoStore) GetAll() ([]*User, error) {
	users := []*User{}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	err := coll.Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//GetByID returns the user with the given ID
func (ms *MongoStore) GetByID(id interface{}) (*User, error) {
	user := &User{}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	err := coll.FindId(id).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	user := &User{}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	err := coll.Find(bson.M{"email": email}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByUserName returns the User with the given username
func (ms *MongoStore) GetByUserName(username string) (*User, error) {
	user := &User{}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	err := coll.Find(bson.M{"username": username}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Insert inserts a new User into the database
//and returns a User with new ID, or an error
func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	user, err := newUser.ToUser()
	if err != nil {
		return nil, err
	}
	user.ID = bson.NewObjectId()
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	if err := coll.Insert(user); err != nil {
		return nil, err
	}
	return user, nil
}

//Update applies UserUpdates to the currentUser
func (ms *MongoStore) Update(updates *UserUpdates, currentuser *User) error {
	if sID, ok := currentuser.ID.(string); ok {
		currentuser.ID = bson.ObjectIdHex(sID)
	}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	if err := coll.UpdateId(currentuser.ID, bson.M{"$set": updates}); err != nil {
		return err
	}
	return nil
}

//ResetPassword applies password resets to the user with the given email
func (ms *MongoStore) ResetPassword(email, newPassword string) error {
	user, err := ms.GetByEmail(email)
	if err != nil {
		return err
	}
	user.SetPassword(newPassword)
	reset := bson.M{"passHash": user.PassHash}
	coll := ms.Session.DB(ms.Database).C(ms.Collection)
	if err := coll.UpdateId(user.ID, bson.M{"$set": reset}); err != nil {
		return err
	}
	return nil
}
