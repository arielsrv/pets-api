package model

type RepositoryDto struct {
	Name       string `json:"name,omitempty"`
	GroupID    int64  `json:"group_id,omitempty"`
	TemplateID int64  `json:"template_id,omitempty"`
}
