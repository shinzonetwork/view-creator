package models

type Metadata struct {
	Version   int        `json:"_v"`
	Total     int        `json:"_t"`
	Revisions []Revision `json:"revisions"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}
