package response

type User struct {
	UserID   string `json:"userID" extensions:"x-order=0"`
	Username string `json:"username" extensions:"x-order=1"`
	Email    string `json:"email" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Role     string `json:"role" extensions:"x-order=4"`
	IsActive bool   `json:"isActive" extensions:"x-order=5"`
	// RegisteredOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	RegisteredOn   string `json:"registeredOn" extensions:"x-order=6"`
	TotalThread    uint   `json:"totalThread" extensions:"x-order=7"`
	TotalFollower  uint   `json:"totalFollower" extensions:"x-order=8"`
	TotalFollowing uint   `json:"totalFollowing" extensions:"x-order=9"`
	IsFollowed     bool   `json:"isFollowed" extensions:"x-order=10"`
}
