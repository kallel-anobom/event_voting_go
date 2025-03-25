package model

type Votes struct {
	ID		string    `json:"id"`
	Name	string `json:"name"`
	Votes	 map[string]int     `json:"votes"`
	EventName	string    `json:"event_name"`
	Date	string `json:"date"`
	Time	string `json:"time"`
}

type VoteSummary struct {
	TotalVotes      int            `json:"total_votes"`
	VotesByOption   map[string]int `json:"votes_by_option"`
	VotesByHour     map[string]int `json:"votes_by_hour"`
	VotesByParticipant map[string]int `json:"votes_by_participant"`
}