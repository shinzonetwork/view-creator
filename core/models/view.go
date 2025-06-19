package models

type View struct {
	Name      string    `json:"name"`
	Query     *string   `json:"query"`
	Sdl       *string   `json:"sdl"`
	Transform Transform `json:"transform"`
	Metadata  Metadata  `json:"metadata"`
}
