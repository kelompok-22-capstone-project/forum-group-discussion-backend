package response

type Moderator struct {
	ID       string `json:"moderatorID" extensions:"x-order=0"`
	UserID   string `json:"userId" extensions:"x-order=1"`
	Username string `json:"username" extensions:"x-order=2"`
	Email    string `json:"email" extensions:"x-order=3"`
	Name     string `json:"name" extensions:"x-order=4"`
	Role     string `json:"role" extensions:"x-order=5"`
	IsActive bool   `json:"isActive" extensions:"x-order=6"`
	// RegisteredOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	RegisteredOn string `json:"registeredOn" extensions:"x-order=7"`
}
