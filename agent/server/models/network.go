package models

type NetworkProcessStat struct {
	PID           int32  `json:"pid"`
	Process       string `json:"process"`
	Username      string `json:"username"`
	Connections   int    `json:"connections"`
	Listening     int    `json:"listening"`
	Established   int    `json:"established"`
	OtherStates   int    `json:"otherStates"`
}


