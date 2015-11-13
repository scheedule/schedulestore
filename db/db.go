// Package db handles all schedule storing and retrieval from the database.
// This package provides an abstraction to allow users to interact with the
// database with the Schedule type and restrict usage to looking up, putting,
// and purging.
package db

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/scheedule/schedulestore/types"
)

var (
	// ScheduleNotFound is returned when a schedule can't be resolved.
	ScheduleNotFound error = errors.New("Schedule Not Found")
)

// Main primitive to hold db connection and attributes. Users will obtain
// and make requests with the DB type.
type DB struct {
	session        *mgo.Session
	collection     *mgo.Collection
	server         string
	dbName         string
	collectionName string
}

// Construct a new DB type
func New(ip, port, dbName, collectionName string) *DB {
	return &DB{
		server:         ip + ":" + port,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

// Initialize connection to the database. An error will be returned if a datbase
// can't be connected to within a minute
func (db *DB) Init() error {
	// Initiate DB connection
	session, err := mgo.DialWithTimeout(db.server, 5*time.Second)
	if err != nil {
		log.Error("failed to dial database:", err)
		return err
	}

	db.session = session

	// Establish Session
	db.collection = db.session.DB(db.dbName).C(db.collectionName)

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

// Delete Schedule from the database
func (db *DB) Delete(userID, name string) error {
	return db.collection.Update(bson.M{
		"user_id": userID,
	}, bson.M{
		"$pull": bson.M{
			"schedules": bson.M{
				"name": name,
			},
		},
	})
}

// Put Schedule into the database
func (db *DB) Put(userID string, entry types.Schedule) error {
	log.Debug("inputing schedule into database: ", entry)

	// Delete any schedule of such name
	err := db.Delete(userID, entry.Name)
	if err != nil {
		log.Warn("delete failed for put.. who cares.")
	}

	_, err = db.collection.Upsert(bson.M{
		"user_id": userID,
	}, bson.M{
		"$addToSet": bson.M{
			"schedules": entry,
		},
	})

	if err != nil {
		log.Error("schedule failed to be input")
	}
	return err
}

// Lookup Schedule in the database
func (db *DB) Lookup(userID string) ([]types.Schedule, error) {
	temp := &types.ScheduleSet{}

	err := db.collection.Find(bson.M{
		"user_id": userID,
	}).One(temp)

	if err != nil {
		log.Debug("did not find schedule for user_id:", userID)
		return nil, ScheduleNotFound
	}

	return temp.Schedules, nil
}
