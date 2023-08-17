package g

import "github.com/khaosles/giz/structs"

/*
   @File: base_service.go
   @Author: khaosles
   @Time: 2023/6/12 01:29
   @Desc: service结构体继承该结构体
*/

type BaseService[T any] struct {
	Dao Dao[T]
}

func (srv BaseService[T]) Save(entity *T) error {
	return srv.Dao.Save(entity)
}

func (srv BaseService[T]) Saves(entities []*T) error {
	return srv.Dao.InsertList(entities)
}

func (srv BaseService[T]) DeleteById(id string) error {
	return srv.Dao.DeleteById(id)
}

func (srv BaseService[T]) DeleteByIds(ids ...string) error {
	return srv.Dao.DeleteByIds(ids...)
}

func (srv BaseService[T]) Update(entity *T) error {
	return srv.Dao.Update(entity)
}

func (srv BaseService[T]) FindById(id string) (*T, error) {
	return srv.Dao.SelectById(id)
}

func (srv BaseService[T]) FindBy(colName string, value any) (*T, error) {
	var record T
	structs.SetField(&record, colName, value)
	return srv.Dao.SelectOne(&record)
}

func (srv BaseService[T]) FindByIds(ids ...string) ([]*T, error) {
	return srv.Dao.SelectByIds(ids...)
}

func (srv BaseService[T]) FindByCondition(conditions *Conditions) ([]*T, error) {
	return srv.Dao.SelectByCondition(conditions)
}

func (srv BaseService[T]) FindAll() ([]*T, error) {
	return srv.Dao.SelectAll()
}
