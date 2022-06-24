package entity

type Entity interface {
	Thread | User | Comment | Category | UserBanned
}

type Pagination[T Entity] struct {
	List     []T
	PageInfo PageInfo
}
