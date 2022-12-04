package model

import "github.com/src/main/app/ent/property"

type App struct {
	ID        int64            `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	URL       string           `json:"url,omitempty"`
	GroupID   int64            `json:"group_id,omitempty"`
	AppTypeID property.AppType `json:"app_type_id,omitempty"`
}
