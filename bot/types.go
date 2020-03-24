package bot

type PollEvent struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Count  string `json:"count"`
}
