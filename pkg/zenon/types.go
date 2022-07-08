package zenon

import "github.com/zenon-network/go-zenon/common/types"

type Event struct {
	MessageType string `json:"type"`
	Id          string `json:"id"`
}

type Votes struct {
	Yes     int `json:"yes"`
	No      int `json:"no"`
	Abstain int `json:"abstain"`
}

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

type ProjectNew struct {
	Id   string  `json:"id"`
	Data Project `json:"data"`
}

type ProjectStatusUpdate struct {
	Id        types.Hash `json:"id"`
	OldStatus uint8      `json:"old"`
	NewStatus uint8      `json:"new"`
}

type PhaseStatusUpdate struct {
	Id        types.Hash `json:"id"`
	Pid       types.Hash `json:"pid"`
	OldStatus uint8      `json:"old"`
	NewStatus uint8      `json:"new"`
}
