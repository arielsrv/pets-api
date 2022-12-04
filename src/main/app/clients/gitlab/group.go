package gitlab

type GroupResponse struct {
	ID   int64  `json:"id"`
	Path string `json:"full_path"`
}
