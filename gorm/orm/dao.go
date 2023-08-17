package g

/*
   @File: dao.go
   @Author: khaosles
   @Time: 2023/6/11 17:20
   @Desc: Dao 接口继承该接口
*/

var _ Dao[any] = (*BaseDao[any])(nil)

type Dao[T any] interface {
	Save(record *T) error                      // 插入或者更新
	Insert(record *T) error                    // 插入
	InsertList(records []*T) error             // 批量插入
	InsertBatch(records []*T, batch int) error // 批量插入, 自定义batch
	InsertOrSelect(record *T) error            // 插入或者获取第一个

	Delete(record *T) error                             // 删除
	DeleteHard(record *T) error                         // 硬删除
	DeleteById(id string) error                         // 根据id删除
	DeleteHardById(id string) error                     // 根据id硬删除
	DeleteByIds(ids ...string) error                    // 根据多个id批量删除
	DeleteHardByIds(ids ...string) error                // 根据多个id批量硬删除
	DeleteByCondition(conditions *Conditions) error     // 根据条件删除
	DeleteHardByCondition(conditions *Conditions) error // 根据条件硬删除

	Update(*T) error                                                    // 更新
	UpdateSelective(record *T, values any) error                        // 更新一个数据多列
	UpdateSelectiveByCondition(record *T, conditions *Conditions) error // 根据条件更新部分字段
	UpdateByCondition(record *T, conditions *Conditions) error          // 根据条件更新全部字段

	SelectById(id string) (*T, error)                                    // 根据id查询
	SelectByIds(ids ...string) ([]*T, error)                             // 根据多个id查询
	SelectOne(record *T) (*T, error)                                     // 根据记录查询一个
	SelectOneByConditions(record *T, conditions *Conditions) (*T, error) // 根据条件查询一个
	Select(record *T) ([]*T, error)                                      // 根据记录查询多个
	SelectAll() ([]*T, error)                                            // 查询所有
	SelectCount(record *T) (int64, error)                                // 查询数量
	SelectByCondition(conditions *Conditions) ([]*T, error)              // 根据条件查询
	SelectCountByCondition(conditions *Conditions) (int64, error)        // 根据条件查询数量
	SelectDistinct(conditions *Conditions) ([]*T, error)                 // 去重查询

	SelectPage(currentPage, pageSize int, sort string) (*Pagination[T], error)                                    // 分页查询
	SelectPageByCondition(currentPage, pageSize int, sort string, conditions *Conditions) (*Pagination[T], error) // 根据条件分页查询

	Exist(record *T) (bool, error) // 判断是否存在
}
