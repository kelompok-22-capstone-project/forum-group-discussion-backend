package payload

type CreateThread struct {
	Title       string `json:"title" validate:"nonzero,min=2,max=50" extensions:"x-order=0"`
	Description string `json:"description" validate:"nonzero,min=2" extensions:"x-order=1"`
	CategoryID  string `json:"categoryID" validate:"nonzero,min=4,max=6" extensions:"x-order=2"`
}
