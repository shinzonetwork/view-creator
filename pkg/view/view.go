package view

import "time"

type Lens struct {
	Path      string         `json:"path"`
	Arguments map[string]any `json:"arguments"`
}

type View struct {
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	QueryFile     string    `json:"queryFile,omitempty"`
	TypeFile      string    `json:"typeFile,omitempty"`
	TransformFile string    `json:"transformFile,omitempty"`
}
