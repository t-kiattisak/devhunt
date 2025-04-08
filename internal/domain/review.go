package domain

type Review struct {
	ID        int    `json:"id"`
	ToolID    int    `json:"tool_id"`
	UserID    string `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}
