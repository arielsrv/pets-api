package model

import "github.com/src/main/app/ent/property"

type CreateAppRequest struct {
	Name      string           `json:"name,omitempty"`
	GroupID   int64            `json:"group_id,omitempty"`
	AppTypeID property.AppType `json:"app_type_id,omitempty"`
}

type CreateAppResponse struct {
	ID  int64  `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

type AppResponse struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
