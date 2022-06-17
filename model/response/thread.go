package response

type ManyThread struct {
	ID           string `json:"ID" extensions:"x-order=0"`
	Title        string `json:"title" extensions:"x-order=1"`
	CategoryID   string `json:"categoryID" extensions:"x-order=2"`
	CategoryName string `json:"categoryName" extensions:"x-order=3"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn     string `json:"publishedOn" extensions:"x-order=4"`
	IsLiked         bool   `json:"isLiked" extensions:"x-order=5"`
	IsFollowed      bool   `json:"isFollowed" extensions:"x-order=6"`
	Description     string `json:"description" extensions:"x-order=7"`
	TotalViewer     uint64 `json:"totalViewer" extensions:"x-order=8"`
	TotalLike       uint64 `json:"totalLike" extensions:"x-order=9"`
	TotalFollower   uint64 `json:"totalFollower" extensions:"x-order=10"`
	TotalComment    uint64 `json:"totalComment" extensions:"x-order=11"`
	CreatorID       string `json:"creatorID" extensions:"x-order=12"`
	CreatorUsername string `json:"creatorUsername" extensions:"x-order=13"`
	CreatorName     string `json:"creatorName" extensions:"x-order=14"`
}

type Thread struct {
	ID           string `json:"ID" extensions:"x-order=0"`
	Title        string `json:"title" extensions:"x-order=1"`
	CategoryID   string `json:"categoryID" extensions:"x-order=2"`
	CategoryName string `json:"categoryName" extensions:"x-order=3"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn     string      `json:"publishedOn" extensions:"x-order=4"`
	IsLiked         bool        `json:"isLiked" extensions:"x-order=5"`
	IsFollowed      bool        `json:"isFollowed" extensions:"x-order=6"`
	Moderators      []Moderator `json:"moderators" extensions:"x-order=7"`
	Description     string      `json:"description" extensions:"x-order=8"`
	TotalViewer     uint64      `json:"totalViewer" extensions:"x-order=9"`
	TotalLike       uint64      `json:"totalLike" extensions:"x-order=10"`
	TotalFollower   uint64      `json:"totalFollower" extensions:"x-order=11"`
	TotalComment    uint64      `json:"totalComment" extensions:"x-order=12"`
	CreatorID       string      `json:"creatorID" extensions:"x-order=13"`
	CreatorUsername string      `json:"creatorUsername" extensions:"x-order=14"`
	CreatorName     string      `json:"creatorName" extensions:"x-order=15"`
}
