package responses

type ProjectResponse struct {
	ID  int64  `json:"id,omitempty"`
	URL string `json:"http_url_to_repo,omitempty"`
}
