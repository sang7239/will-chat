package handlers

import (
	"github.com/will-slack/apiserver/events"
	"github.com/will-slack/apiserver/models/messages"
	"github.com/will-slack/apiserver/models/users"
	"github.com/will-slack/apiserver/sessions"
)

// Context holds a store for the server
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	MessageStore messages.Store
	Notifier     *events.Notifier
}
