package model

// AppModel Model
// swagger:model AppModel
type AppModel struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
