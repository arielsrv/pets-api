package request

import "github.com/src/main/app/ent/property"

type CreateAppRequest struct {
	Name      string           `json:"name,omitempty"`
	GroupID   int64            `json:"group_id,omitempty"`
	AppTypeID property.AppType `json:"app_type_id,omitempty"`
}
