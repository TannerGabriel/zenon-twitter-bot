package zenon

import "github.com/zenon-network/go-zenon/common/types"

// Event is the data type for any ZMQ message
type Event struct {
	MessageType string `json:"type"`
	Id          string `json:"id"`
}

// Votes contains the fields of a AZ vote
type Votes struct {
	Yes     int `json:"yes"`
	No      int `json:"no"`
	Abstain int `json:"abstain"`
}

// Project defines an AZ project
type Project struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	Created     int      `json:"created"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Znn         float64  `json:"znn"`
	Qsr         float64  `json:"qsr"`
	Status      int      `json:"status"`
	Votes       Votes    `json:"votes"`
	Phases      []string `json:"phases"`
}

// ProjectNew contains the event data of a new project
type ProjectNew struct {
	Id   string  `json:"id"`
	Data Project `json:"data"`
}

// ProjectStatusUpdate contains the fields of the project:status-update event
type ProjectStatusUpdate struct {
	Id        types.Hash `json:"id"`
	OldStatus uint8      `json:"old"`
	NewStatus uint8      `json:"new"`
}

// PhaseStatusUpdate contains the fields of the phase:status-update event
type PhaseStatusUpdate struct {
	Id        types.Hash `json:"id"`
	Pid       types.Hash `json:"pid"`
	OldStatus uint8      `json:"old"`
	NewStatus uint8      `json:"new"`
}
