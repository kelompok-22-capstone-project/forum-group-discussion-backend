package response

type PageInfo struct {
	Limit     uint `json:"limit" extensions:"x-order=0"`
	Page      uint `json:"page" extensions:"x-order=1"`
	PageTotal uint `json:"pageTotal" extensions:"x-order=2"`
	Total     uint `json:"total" extensions:"x-order=3"`
}
