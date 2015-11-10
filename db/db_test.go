package db

import (
	"github.com/scheedule/schedulestore/types"
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	mydb := NewDB(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	if mydb == nil {
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	mydb := NewDB(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	err := mydb.Init()
	if err != nil {
		t.Error("Failed to initialize DB:", err)
	}
	if mydb.session == nil {
		t.Error("Session is nil")
	}
	if mydb.collection == nil {
		t.Error("Collection is nil")
	}
}

func getDB() *DB {
	mydb := NewDB(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	_ = mydb.Init()
	return mydb
}

var sampleSchedule = types.Schedule{
	UserID: "1",
	CRNs:   []string{"1", "2"},
}

func TestPurge(t *testing.T) {
	mydb := getDB()
	mydb.Put(sampleSchedule)
	mydb.Purge()
	_, err := mydb.Lookup("1")
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

	mydb := getDB()
	mydb.Close()
	_ = mydb.session.Ping()
}

func TestPut(t *testing.T) {
	mydb := getDB()
	mydb.Purge()

	err := mydb.Put(sampleSchedule)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
}

func TestLookup(t *testing.T) {
	mydb := getDB()
	mydb.Purge()

	err := mydb.Put(sampleSchedule)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
	schedule, err := mydb.Lookup("1")
	if err != nil {
		t.Error("Schedule lookup returned error: ", err)
	}
	if schedule.CRNs[0] != "1" || schedule.CRNs[1] != "2" {
		t.Error("Lookup result inaccurate: ", schedule)
	}
}
