package model

import "github.com/src/main/ent/property"

type CreateAppModel struct {
	Name      string           `json:"name,omitempty"`
	GroupID   int64            `json:"group_id,omitempty"`
	AppTypeID property.AppType `json:"app_type_id,omitempty"`
}
