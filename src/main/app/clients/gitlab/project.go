package gitlab

type CreateProjectRequest struct {
	Name    string `json:"name,omitempty"`
	GroupID int64  `json:"namespace_id,omitempty"`
}

type CreateProjectResponse struct {
	ID  int64  `json:"id,omitempty"`
	URL string `json:"http_url_to_repo,omitempty"`
}

type ProjectResponse struct {
	ID  int64  `json:"id,omitempty"`
	URL string `json:"http_url_to_repo,omitempty"`
}
