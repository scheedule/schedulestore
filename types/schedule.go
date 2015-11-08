// Package types holds types shared across the schedule store. Types are tagged
// for JSON unmarshalling and bson serializing.
package types

type (
	// JSON query type to unmarshal.
	ScheduleProposal []string

	// Schedule type to be serialized to BSON and placed in db.
	Schedule struct {
		UserID string   `bson:"user_id"`
		CRNs   []string `bson:"crns"`
	}
)
