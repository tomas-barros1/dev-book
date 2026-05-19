package dtos

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"max=2000"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=200"`
	Content string `json:"content" binding:"omitempty,max=2000"`
}

type PostResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  uint   `json:"author_id"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}
