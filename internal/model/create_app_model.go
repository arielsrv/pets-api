package model

type CreateAppModel struct {
	Name      string `json:"name,omitempty"`
	GroupID   int64  `json:"group_id,omitempty"`
	AppTypeID int    `json:"app_type_id,omitempty"`
}
