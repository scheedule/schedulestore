// Package API offers routes to expose an API for users to store and retrieve
// schedules. The "user_id" header must be set to know who's schedule is being
// retrieved or saved.
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/scheedule/schedulestore/db"
	"github.com/scheedule/schedulestore/types"
)

var (
	BadRequestError   = errors.New("The request was malformed.")
	DBError           = errors.New("Query to database failed.")
	UnmarshalError    = errors.New("Error unmarshalling data from the database.")
	UnauthorizedError = errors.New("Unauthorized.")
	errorMap          = map[error]int{
		UnauthorizedError: http.StatusUnauthorized,
		BadRequestError:   http.StatusBadRequest,
		DBError:           http.StatusNotFound,
		UnmarshalError:    http.StatusInternalServerError,
	}
)

// Simply contains a DB. Mainly for attaching functions.
type API struct {
	db *db.DB
}

func New(db *db.DB) *API {
	return &API{db}
}

func (a *API) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(r)
	log.Debug("received request: ", r)
	if err != nil {
		handleError(w, UnauthorizedError)
		return
	}

	switch r.Method {
	case "GET":
		a.handleGET(userID, w, r)
	case "PUT":
		a.handlePUT(userID, w, r)
	case "DELETE":
		a.handleDELETE(userID, w, r)
	default:
		handleError(w, BadRequestError)
	}
}

func extractUserID(r *http.Request) (string, error) {

	userID := r.Header.Get("user_id")
	if userID == "" {
		log.Warn("user_id not set, rejecting.")
		return "", BadRequestError
	}

	return userID, nil
}

func writeJSON(w http.ResponseWriter, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

// HTTP handler to lookup schedule and write it as a response
func (a *API) handleGET(userID string, w http.ResponseWriter, r *http.Request) {

	schedules, err := a.db.Lookup(userID)
	if err != nil {
		handleError(w, err)
		return
	}

	js, err := json.Marshal(schedules)
	if err != nil {
		handleError(w, UnmarshalError)
		return
	}

	writeJSON(w, js)
}

// HTTP handler to allow clients to add/update schedules
func (a *API) handlePUT(userID string, w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	log.Debug("going to decode schedule")
	decoder := json.NewDecoder(r.Body)
	proposal := types.Schedule{}
	err := decoder.Decode(&proposal)
	if err != nil {
		log.Warn("failed to decode proposed schedule: ", err)
		handleError(w, UnmarshalError)
		return
	}

	log.Debug("proposal: ", proposal)
	err = a.db.Put(userID, proposal)
	if err != nil {
		log.Warn("error putting: ", err)
		handleError(w, err)
		return
	}

	w.WriteHeader(200)
}

// HHTP handler to allow clients to delete schedules
func (a *API) handleDELETE(userID string, w http.ResponseWriter, r *http.Request) {
	scheduleName := r.FormValue("name")
	if scheduleName == "" {
		handleError(w, BadRequestError)
		return
	}

	err := a.db.Delete(userID, scheduleName)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(200)
}

// Write the appropriate message to the client
func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), errorMap[err])
}
