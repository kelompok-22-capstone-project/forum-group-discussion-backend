package payload

type CreateComment struct {
	Comment string `json:"comment" validate:"nonzero,min=1" extensions:"x-order=0"`
}
