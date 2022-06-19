package response

type Entity interface {
	ManyThread | Category | Comment | User
}

type Pagination[T Entity] struct {
	List     []T      `json:"list" extensions:"x-order=0"`
	PageInfo PageInfo `json:"pageInfo" extensions:"x-order=1"`
}
