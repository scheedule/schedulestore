package db

import (
	"github.com/scheedule/schedulestore/types"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	if myDB == nil {
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	err := myDB.Init()
	if err != nil {
		t.Error("Failed to initialize DB:", err)
	}
	if myDB.session == nil {
		t.Error("Session is nil")
	}
	if myDB.collection == nil {
		t.Error("Collection is nil")
	}
}

func getDB() *DB {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	_ = myDB.Init()
	return myDB
}

var sampleSchedule = types.Schedule{
	Name:    "1",
	CRNList: []string{"1", "2"},
}

func TestPurge(t *testing.T) {
	myDB := getDB()
	myDB.Put("fake", sampleSchedule)
	myDB.Purge()
	_, err := myDB.Lookup("1")
	if err == nil {
		t.Errorf("Database should be empty but returned no error when looked up schedule.")
	}
}

func TestClose(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Using a closed session should panic")
		}
	}()

	myDB := getDB()
	myDB.Close()
	_ = myDB.session.Ping()
}

func TestPut(t *testing.T) {
	myDB := getDB()
	myDB.Purge()

	err := myDB.Put("fake", sampleSchedule)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
}

func TestLookup(t *testing.T) {
	myDB := getDB()
	myDB.Purge()

	err := myDB.Put("fake", sampleSchedule)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
	_, err = myDB.Lookup("fake")
	if err != nil {
		t.Error("Schedule lookup returned error: ", err)
	}
}
