package payload

type User struct {
	Username string `json:"username" validate:"nonzero,min=2,max=20" extensions:"x-order=0"`
}
