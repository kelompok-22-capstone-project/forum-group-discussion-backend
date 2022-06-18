package payload

type CreateComment struct {
	Comment string `json:"name" validate:"nonzero,min=1" extensions:"x-order=0"`
}
