package response

type Category struct {
	ID          string `json:"ID" extensions:"x-order=0"`
	Name        string `json:"name" extensions:"x-order=1"`
	Description string `json:"description" extensions:"x-order=2"`
	// CreatedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	CreatedOn string `json:"createdOn" extensions:"x-order=3"`
}
