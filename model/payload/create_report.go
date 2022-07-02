package payload

type CreateReport struct {
	Username  string `json:"username" validate:"nonzero,min=2,max=10" extensions:"x-order=0"`
	CommentID string `json:"commentID" validate:"nonzero,min=2,max=10" extensions:"x-order=1"`
	Reason    string `json:"reason" validate:"nonzero,min=2" extensions:"x-order=2"`
}
