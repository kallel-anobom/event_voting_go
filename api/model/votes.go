package model

import "time"

type Vote struct {
	ParticipantID int `json:"participant_id"`
	Date          time.Time
}
