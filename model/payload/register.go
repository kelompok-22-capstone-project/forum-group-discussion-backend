package payload

type Register struct {
	Username string `json:"username" validate:"nonzero,min=2,max=20" extensions:"x-order=0"`
	Email    string `json:"email" validate:"nonzero,min=5,max=50" extensions:"x-order=1"`
	Name     string `json:"name" validate:"nonzero,min=1,max=50" extensions:"x-order=2"`
	Password string `json:"password" validate:"nonzero,min=5,max=20" extensions:"x-order=3"`
}
