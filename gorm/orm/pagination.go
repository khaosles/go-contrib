package g

/*
   @File: pagination.go
   @Author: khaosles
   @Time: 2023/4/30 09:07
   @Desc: 分页器，datacollection 为返回数据
*/

type Pagination[T any] struct {
	PageSize       int   `json:"pageSize"`       // 每页大小
	CurrentPage    int   `json:"currentPage"`    // 当前页
	TotalCount     int64 `json:"totalCount"`     // 总数
	TotalPages     int64 `json:"totalPages"`     // 总页数
	DataCollection []*T  `json:"dataCollection"` // 当前页数据
}
