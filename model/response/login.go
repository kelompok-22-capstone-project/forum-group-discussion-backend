package response

type Login struct {
	Token string `json:"token" extensions:"x-order=0"`
	Role  string `json:"role" extensions:"x-order=1"`
}
