package payload

type CreateCategory struct {
	Name        string `json:"name" validate:"nonzero,min=2,max=50" extensions:"x-order=0"`
	Description string `json:"description" validate:"nonzero,min=2" extensions:"x-order=1"`
}
