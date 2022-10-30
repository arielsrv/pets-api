package model

type CreateProjectModel struct {
	ID  int64  `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}
