package model

type GroupDto struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	NamespaceID int64  `json:"namespace_id,omitempty"`
}
