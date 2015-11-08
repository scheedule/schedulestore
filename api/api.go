// Package api offers routes to expose an API for users to store and retrieve
// schedules. The "user_id" header must be set to know who's schedule is being
// retrieved or saved.
package api

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/schedulestore/db"
	"github.com/scheedule/schedulestore/types"
	"io/ioutil"
	"net/http"
)

var (
	BadRequestError = errors.New("The request was malformed")
	DBError         = errors.New("Query to database failed.")
	UnmarshalError  = errors.New("Error unmarshalling data from the database")
)

// Write the appropriate message to the client
func handleError(w http.ResponseWriter, err error) {
	switch err {
	case BadRequestError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case DBError:
		http.Error(w, err.Error(), http.StatusNotFound)
	case UnmarshalError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Simply contains a DB. Mainly for attaching functions.
type Api struct {
	Mydb *db.DB
}

// HTTP Handler to lookup schedule and write it as a response
func (a *Api) HandleLookup(w http.ResponseWriter, r *http.Request) {
	user_id := r.Header.Get("user_id")
	if user_id == "" {
		handleError(w, BadRequestError)
		return
	}
	log.Debug("Lookup from user_id:", user_id)

	js, err := lookupSchedule(a.Mydb, user_id)
	if err != nil {
		log.Warn("Lookup failed:", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	log.Debug("Lookup successful")
}

// HTTP Handler to put schedule and respond with success indication
func (a *Api) HandlePut(w http.ResponseWriter, r *http.Request) {
	user_id := r.Header.Get("user_id")
	if user_id == "" {
		handleError(w, BadRequestError)
		return
	}
	log.Debug("Put from user_id:", user_id)

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn("Failed to read request:", err)
		handleError(w, BadRequestError)
		return
	}

	proposal := types.ScheduleProposal{}
	err = json.Unmarshal([]byte(body), &proposal)
	if err != nil {
		log.Warn("Failed to unmarshal request:", err)
		handleError(w, UnmarshalError)
		return
	}

	err = putSchedule(a.Mydb, user_id, proposal)
	if err != nil {
		log.Warn("Failed to put schedule:", err)
		handleError(w, err)
		return
	}

	w.WriteHeader(200)
	log.Debug("Schedule successfully put in DB")
}

// Helper function to lookup schedule and Marshal given db and user_id
func lookupSchedule(db *db.DB, user_id string) ([]byte, error) {
	schedule, err := db.Lookup(user_id)
	if err != nil {
		return nil, DBError
	}

	js, err := json.Marshal(schedule.CRNs)
	if err != nil {
		return nil, UnmarshalError
	}

	return js, nil
}

// Helper function to put schedule into the db
func putSchedule(db *db.DB, user_id string, schedule []string) error {
	sch := types.Schedule{
		UserID: user_id,
		CRNs:   schedule,
	}

	if err := db.Put(sch); err != nil {
		return DBError
	}
	return nil
}
