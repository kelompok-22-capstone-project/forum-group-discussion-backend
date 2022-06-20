package payload

type CreateReport struct {
	UserID   string `json:"userID" validate:"nonzero,min=2,max=10" extensions:"x-order=0"`
	ThreadID string `json:"threadID" validate:"nonzero,min=2,max=10" extensions:"x-order=1"`
	Reason   string `json:"reason" validate:"nonzero,min=2" extensions:"x-order=2"`
}
