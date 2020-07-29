package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/will-slack/apiserver/models/users"
	"github.com/will-slack/apiserver/sessions"
)

// UserHandler allows new users to sign up (POST) or return all users (GET)
func (ctx *Context) UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		newUser := &users.NewUser{}
		if err := decoder.Decode(newUser); err != nil {
			http.Error(w, "Unable to decode response", http.StatusBadRequest)
			return
		}
		if err := newUser.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if u, _ := ctx.UserStore.GetByEmail(newUser.Email); u != nil {
			http.Error(w, "Invalid Email. Already exists", http.StatusBadRequest)
			return
		}
		if u, _ := ctx.UserStore.GetByUserName(newUser.UserName); u != nil {
			http.Error(w, "Invalid username. Already exists", http.StatusBadRequest)
			return
		}
		var user *users.User
		var err error
		if user, err = ctx.UserStore.Insert(newUser); err != nil {
			http.Error(w, "Unable to register new user", http.StatusInternalServerError)
			return
		}
		state := &SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}
		if _, err := sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w); err != nil {
			http.Error(w, "Unable to begin session for the user", http.StatusInternalServerError)
			return
		}
		Respond(w, user, "application/json; charset=utf-8")

	} else if r.Method == "GET" {
		var users []*users.User
		var err error
		if users, err = ctx.UserStore.GetAll(); err != nil {
			http.Error(w, "Unable to get all users", http.StatusInternalServerError)
			return
		}
		Respond(w, users, "application/json; charset=utf-8")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

// SpecificUserHandler maps given username with its id
func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	_, username := path.Split(r.URL.Path)
	if r.Method == "GET" {
		var user *users.User
		var err error
		if user, err = ctx.UserStore.GetByUserName(username); err != nil {
			http.Error(w, "Unable to find given user", http.StatusBadRequest)
			return
		}
		Respond(w, user.ID, "application/json; charset=utf-8")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

// SessionHandler allows existing users to sign-in
func (ctx *Context) SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		credentials := &users.Credentials{}
		if err := decoder.Decode(credentials); err != nil {
			http.Error(w, "Unable to decode response", http.StatusBadRequest)
			return
		}
		var user *users.User
		var err error
		if user, err = ctx.UserStore.GetByEmail(credentials.Email); err != nil {
			http.Error(w, "User Not Authorized", http.StatusUnauthorized)
			return
		}
		err = user.Authenticate(credentials.Password)
		if err != nil {
			http.Error(w, "User Not Authorized", http.StatusUnauthorized)
			return
		}
		state := &SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}
		if _, err := sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w); err != nil {
			http.Error(w, "Unable to begin session for the user", http.StatusInternalServerError)
			return
		}
		Respond(w, user, "application/json; charset=utf-8")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

// SessionMineHandler allows authenticated users to sign-out
func (ctx *Context) SessionMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		if _, err := sessions.EndSession(r, ctx.SessionKey, ctx.SessionStore); err != nil {
			http.Error(w, "Unable to end session for the user", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "signed out")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

// UserMeHandler gets the session state and responds with session state's user field
func (ctx *Context) UserMeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		state := &SessionState{}
		if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		Respond(w, state.User, "application/json; charset=utf-8")
	} else if r.Method == "PATCH" {
		state := &SessionState{}
		sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		decoder := json.NewDecoder(r.Body)
		updates := &users.UserUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "Error while decoding", http.StatusBadRequest)
			return
		}
		if err := ctx.UserStore.Update(updates, state.User); err != nil {
			http.Error(w, "Error in db while updating user", http.StatusInternalServerError)
			return
		}
		state.User.FirstName = updates.FirstName
		state.User.LastName = updates.LastName
		if err := ctx.SessionStore.Save(sid, state); err != nil {
			http.Error(w, "Error in redis while updating user", http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "user successfully updated")
	}
}
