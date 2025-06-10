package dto

type TaskStatusResp struct {
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Duration  string `json:"duration"`
}
