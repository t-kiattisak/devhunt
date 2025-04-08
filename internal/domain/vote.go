package domain

type Vote struct {
	ID        int    `json:"id"`
	ToolID    int    `json:"tool_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
