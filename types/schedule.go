// Package types holds types shared across the schedule store. Types are tagged
// for JSON unmarshalling and bson serializing.
package types

type (
	// Schedule type to be serialized to BSON and placed in db.
	Schedule struct {
		Name    string   `bson:"name" json:"name"`
		CRNList []string `bson:"CRNList" json:"CRNList"`
	}

	// ScheduleSet contains an array of schedules.
	ScheduleSet struct {
		UserID    string     `bson:"user_id" json="-"`
		Schedules []Schedule `bson:"schedules" json="schedules"`
	}
)
