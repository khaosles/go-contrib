package g

/*
   @File: service.go
   @Author: khaosles
   @Time: 2023/6/12 01:28
   @Desc:  service 接口继承该接口
*/

type Service[T any] interface {
	Save(entity *T) error                                 // 保存
	Saves(entities []*T) error                            // 批量保存
	DeleteById(id string) error                           // 根据id删除单个
	DeleteByIds(ids ...string) error                      // 根据多个id删除
	Update(entity *T) error                               // 更新
	FindById(id string) (*T, error)                       // 根据id查找
	FindBy(colName string, value any) (*T, error)         // 根据某个字段查找唯一值
	FindByIds(ids ...string) ([]*T, error)                // 根据id查找多个
	FindByCondition(conditions *Conditions) ([]*T, error) // 根据条件查找
	FindAll() ([]*T, error)                               // 查找全部
}
