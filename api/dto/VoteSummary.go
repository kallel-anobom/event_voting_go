package dto

type VoteSummary struct {
	TotalVotes         int            `json:"total_votes"`
	VotesByOption      map[string]int `json:"votes_by_option"`
	VotesByHour        map[string]int `json:"votes_by_hour"`
	VotesByParticipant map[string]int `json:"votes_by_participant"`
}

type VoteMessage struct {
	ParticipantID int `json:"participant_id"`
}
