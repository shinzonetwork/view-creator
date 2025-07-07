package models

type Revision struct {
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
	Diff      string `json:"diff"`
}
