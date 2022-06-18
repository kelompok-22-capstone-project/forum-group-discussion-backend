package response

type Comment struct {
	ID       string `json:"ID" extensions:"x-order=0"`
	UserID   string `json:"userID" extensions:"x-order=1"`
	Username string `json:"username" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Comment  string `json:"comment" extensions:"x-order=4"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn string `json:"publishedOn" extensions:"x-order=5"`
}
