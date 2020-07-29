package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/will-slack/apiserver/models/messages"
	"github.com/will-slack/apiserver/sessions"
)

//ChannelsHandler gets the channels the current user can see and write the returned slice to the
//response as a JSON-encoded array. (GET)
//Add the current user to the new channel's Members list, insert the new channel,
//and write the newly-inserted Channel object to the response (POST)
func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
		http.Error(w, "Unable to find current user", http.StatusForbidden)
		return
	}
	if r.Method == "GET" {
		var channels []*messages.Channel
		var err error
		if channels, err = ctx.MessageStore.GetAll(state.User); err != nil {
			http.Error(w, "Unable to retrieve channels for the user", http.StatusInternalServerError)
		}
		Respond(w, channels, "application/json; charset=utf-8")
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		newChannel := &messages.NewChannel{}
		if err := decoder.Decode(newChannel); err != nil {
			http.Error(w, "Unable to decode response", http.StatusBadRequest)
			return
		}
		if err := newChannel.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newChannel.Members = append(newChannel.Members, state.User.ID)
		var channel *messages.Channel
		var err error
		if channel, err = ctx.MessageStore.InsertChannel(newChannel, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Respond(w, channel, "application/json; charset=utf-8")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

//SpecificChannelHandler updates the specified message if the current user is the message creator (PATCH)
//or deletes the message (DELETE)
func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	_, channelID := path.Split(r.URL.Path)
	state := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
		http.Error(w, "Unable to find current user", http.StatusForbidden)
		return
	}
	if r.Method == "GET" {
		var messages []*messages.Message
		var err error
		if messages, err = ctx.MessageStore.GetMostRecentMessages(state.User, channelID, 500); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		Respond(w, messages, "application/json; charset=utf-8")
	} else if r.Method == "PATCH" {
		decoder := json.NewDecoder(r.Body)
		updates := &messages.ChannelUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "Error while decoding", http.StatusBadRequest)
			return
		}
		if err := ctx.MessageStore.UpdateChannel(updates, channelID, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var channel *messages.Channel
		var err error
		if channel, err = ctx.MessageStore.GetChannelByID(channelID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Respond(w, channel, "application/json; charset=utf-8")
	} else if r.Method == "DELETE" {
		if err := ctx.MessageStore.DeleteChannel(channelID, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		io.WriteString(w, "channel successfully deleted")
	} else if r.Method == "LINK" {
		member := r.Header.Get("Link")
		// add user to public channel
		if len(member) == 0 {
			if err := ctx.MessageStore.AddUserToChannel(state.User.ID, channelID, state.User.ID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else { // add user to private channel
			if err := ctx.MessageStore.AddUserToChannel(member, channelID, state.User.ID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		io.WriteString(w, "User Successfully added to channel")
	} else if r.Method == "UNLINK" {
		member := r.Header.Get("Link")
		// add user to public channel
		if len(member) == 0 {
			if err := ctx.MessageStore.RemoveUserFromChannel(state.User.ID, channelID, state.User.ID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else { // add user to private channel
			if err := ctx.MessageStore.RemoveUserFromChannel(member, channelID, state.User.ID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		io.WriteString(w, "User Successfully removed from channel")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}

// SpecificChannelInfoHandler does this
func (ctx *Context) SpecificChannelInfoHandler(w http.ResponseWriter, r *http.Request) {
	_, channelID := path.Split(r.URL.Path)
	state := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
		http.Error(w, "Unable to find current user", http.StatusForbidden)
		return
	}
	if r.Method == "GET" {
		var channel *messages.Channel
		var err error
		if channel, err = ctx.MessageStore.GetChannelByID(channelID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Respond(w, channel, "application/json; charset=utf-8")
	}
}

//MessagesHandler does this
func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
		http.Error(w, "Unable to find current user", http.StatusForbidden)
		return
	}
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		newMessage := &messages.NewMessage{}
		if err := decoder.Decode(newMessage); err != nil {
			http.Error(w, "Unable to decode response", http.StatusBadRequest)
			return
		}
		if err := newMessage.Validate(); err != nil {
			http.Error(w, "Unable to validate message", http.StatusBadRequest)
			return
		}
		var message *messages.Message
		var err error
		if message, err = ctx.MessageStore.InsertMessage(newMessage, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx.notify("new message", message)
		Respond(w, message, "application/json; charset=utf-8")
	} else {
		http.Error(w, "Invalid Requst", http.StatusBadRequest)
		return
	}
}

//SpecificMessageHandler does this
func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	_, messageID := path.Split(r.URL.Path)
	state := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &state); err != nil {
		http.Error(w, "Unable to find current user", http.StatusForbidden)
		return
	}
	if r.Method == "PATCH" {
		decoder := json.NewDecoder(r.Body)
		updates := &messages.MessageUpdates{}
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "Error while decoding", http.StatusBadRequest)
			return
		}
		if err := ctx.MessageStore.UpdateMessage(updates, messageID, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var message *messages.Message
		var err error
		if message, err = ctx.MessageStore.GetMessageByID(messageID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Respond(w, message, "application/json; charset=utf-8")
	} else if r.Method == "DELETE" {
		if err := ctx.MessageStore.DeleteMessage(messageID, state.User); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		io.WriteString(w, "message successfully deleted")
	} else {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
}
