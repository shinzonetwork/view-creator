package models

type Lens struct {
	Label     string         `json:"label"`
	Path      string         `json:"path"`
	Arguments map[string]any `json:"arguments"`
}
