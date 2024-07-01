package model

type Contest struct {
	Key         string     `json:"_key"`
	ID          string     `json:"_id"`
	Rev         string     `json:"_rev"`
	Name        string     `json:"name"`
	Start       string     `json:"start"`
	Startoffset string     `json:"startoffset"`
	Stop        string     `json:"stop"`
	Stopoffset  string     `json:"stopoffset"`
	Outcomes    []*Outcome `json:"outcomes"`
	Games       []*Game    `json:"games"`
	Venue       *Venue     `json:"venue,omitempty"`
}

type Outcome struct {
	Key    string  `json:"_key"`
	ID     string  `json:"_id"`
	Rev    string  `json:"_rev"`
	Player *Player `json:"player"`
	Place  int     `json:"place"`
	Result string  `json:"result"`
}
