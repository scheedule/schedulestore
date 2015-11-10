// Package db handles all schedule storing and retrieval from the database.
// This package provides an abstraction to allow users to interact with the
// database with the Schedule type and restrict usage to looking up, putting,
// and purging.
package db

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/schedulestore/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	// ScheduleNotFound is returned when a schedule can't be resolved.
	ScheduleNotFound error = errors.New("Schedule Not Found")
)

// Main primitive to hold db connection and attributes. Users will obtain
// and make requests with the DB type.
type DB struct {
	session         *mgo.Session
	collection      *mgo.Collection
	server          string
	db_name         string
	collection_name string
}

// Construct a new DB type
func NewDB(ip, port, db_name, collection_name string) *DB {
	return &DB{
		server:          ip + ":" + port,
		db_name:         db_name,
		collection_name: collection_name,
	}
}

// Initialize connection to the database. An error will be returned if a datbase
// can't be connected to within a minute
func (db *DB) Init() error {
	// Initiate DB connection
	session, err := mgo.DialWithTimeout(db.server, 5*time.Second)
	if err != nil {
		log.Error("Failed to dial database:", err)
		return err
	}

	db.session = session

	// Establish Session
	db.collection = db.session.DB(db.db_name).C(db.collection_name)

	return nil
}

// Drop the specified collection from the database.
func (db *DB) Purge() {
	db.collection.DropCollection()
}

// Close the session with the database
func (db *DB) Close() {
	db.session.Close()
}

// Put Schedule into the database
func (db *DB) Put(entry types.Schedule) error {
	log.Debug("Inputing Schedule into database:", entry)
	_, err := db.collection.Upsert(bson.M{
		"user_id": entry.UserID,
	}, entry)
	if err != nil {
		log.Error("Schedule failed to be input")
	}
	return err
}

// Lookup Schedule in the database
func (db *DB) Lookup(user_id string) (*types.Schedule, error) {
	temp := &types.Schedule{}

	err := db.collection.Find(bson.M{
		"user_id": user_id,
	}).One(temp)

	if err != nil {
		log.Debug("Did not find schedule for user_id:", user_id)
		return nil, ScheduleNotFound
	}

	return temp, nil
}
