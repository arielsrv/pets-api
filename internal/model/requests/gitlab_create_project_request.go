package requests

type CreateProjectRequest struct {
	Name    string `json:"name,omitempty"`
	GroupID int64  `json:"namespace_id,omitempty"`
}
