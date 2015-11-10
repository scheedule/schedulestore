package api

import (
	"github.com/scheedule/schedulestore/db"
	"github.com/scheedule/schedulestore/types"
	"os"
	"testing"
)

var myApi *Api

var sampleSchedule = types.Schedule{
	UserID: "1",
	CRNs:   []string{"1", "2"},
}

func init() {
	mydb := db.NewDB(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	mydb.Init()
	myApi = &Api{mydb}
}

func TestLookupSchedule(t *testing.T) {
	myApi.Mydb.Purge()

	_, err := lookupSchedule(myApi.Mydb, "fakeid")
	if err == nil {
		t.Error("Looking up schedule with fake ID failed to produce an error.")
	}
}

func TestPutSchedule(t *testing.T) {
	myApi.Mydb.Purge()

	err := putSchedule(myApi.Mydb, "1", []string{"1", "2"})
	if err != nil {
		t.Error("Failed to put schedule into database")
	}

	_, err = lookupSchedule(myApi.Mydb, "1")
	if err != nil {
		t.Error("Failure looking up schedule.")
	}
}
