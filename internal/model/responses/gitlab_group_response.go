package responses

type GitLabGroupResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"full_path"`
	Namespace Namespace `json:"namespace"`
}

type Namespace struct {
	ID int64 `json:"id,omitempty"`
}
