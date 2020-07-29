package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("scheme used in Authorization header is not supported")

//BeginSession creates a new session ID, saves the state to the store, adds a
//header to the response with the session ID, and returns the new session ID
func BeginSession(signingKey string, store Store, state interface{}, w http.ResponseWriter) (SessionID, error) {
	//create a new SessionID
	//if you get an error, return InvalidSessionID and the error
	SessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//save the state to the store
	//if you get an error, return InvalidSessionID and the error
	err = store.Save(SessionID, state)
	if err != nil {
		return InvalidSessionID, err
	}
	w.Header().Add(headerAuthorization, fmt.Sprintf("%s%s", schemeBearer, SessionID))
	//Add a response header like this:
	//  Authorization: Bearer <sid>
	//where <sid> is the new SessionID
	//return the new SessionID and nil
	return SessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//get the value of the Authorization header

	auth := r.Header.Get(headerAuthorization)
	//if it's zero-length, check to see if there is a query parameter
	// otherwise return the error
	if len(auth) == 0 {
		query := r.URL.Query().Get("auth")
		if len(query) == 0 {
			return InvalidSessionID, ErrNoSessionID
		}
		// make sure to put the bearer token so the rest of this doesn't break
		auth = schemeBearer + query
	}

	//if it doesn't start with "Bearer ",
	//return InvalidSessionID and ErrInvalidScheme
	if !strings.HasPrefix(auth, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}
	id := strings.TrimPrefix(auth, schemeBearer)
	SessionID, err := ValidateID(id, signingKey)
	//trim off the "Bearer " prefix and validate the remaining id
	//if you get an error return InvalidSessionID and the error
	if err != nil {
		return InvalidSessionID, err
	}
	//return the validated SessionID and nil
	return SessionID, nil
}

//GetState extracts the SessionID from the request,
//and gets the associated state from the provided store
func GetState(r *http.Request, signingKey string, store Store, state interface{}) (SessionID, error) {
	//get the SessionID from the request
	//if you get an error, return the SessionID and error
	SessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return SessionID, err
	}
	//get the associated state data from the provided store
	//if you get an error return the SessionID and the error
	err = store.Get(SessionID, state)
	if err != nil {
		return SessionID, err
	}
	//return the SessionID and nil
	return SessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//get the SessionID from the request
	//if you get an error return the SessionID and error
	SessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return SessionID, err
	}
	//delete the associated data in the provided store
	err = store.Delete(SessionID)
	if err != nil {
		return SessionID, err
	}
	//return the SessionID and nil
	return SessionID, nil

}
