package handlers

import (
	"time"

	"github.com/will-slack/apiserver/models/users"
)

// SessionState holds current session for a user
type SessionState struct {
	BeganAt    time.Time
	ClientAddr string
	User       *users.User
}
