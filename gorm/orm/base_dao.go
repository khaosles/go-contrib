package g

import (
	"github.com/khaosles/giz/structs"
	"gorm.io/gorm"
)

/*
   @File: base_dao.go
   @Author: khaosles
   @Time: 2023/6/11 19:28
   @Desc: BaseDao 结构体继承该结构体
*/

type BaseDao[T any] struct {
	DB *gorm.DB
}

func (dao BaseDao[T]) Save(record *T) error {
	return dao.DB.Save(record).Error
}

func (dao BaseDao[T]) Insert(record *T) error {
	return dao.DB.Create(record).Error
}

func (dao BaseDao[T]) InsertList(records []*T) error {
	return dao.DB.CreateInBatches(records, len(records)).Error
}

func (dao BaseDao[T]) InsertBatch(records []*T, batch int) error {
	if batch < 1 {
		batch = len(records)
	}
	return dao.DB.CreateInBatches(records, batch).Error
}

func (dao BaseDao[T]) InsertOrSelect(record *T) error {
	return dao.DB.FirstOrCreate(record, record).Error
}

func (dao BaseDao[T]) Delete(record *T) error {
	return dao.DB.Delete(record).Error
}

func (dao BaseDao[T]) DeleteHard(record *T) error {
	return dao.DB.Unscoped().Delete(record).Error
}

func (dao BaseDao[T]) DeleteById(id string) error {
	return dao.DB.Delete(new(T), "id = ?", id).Error
}

func (dao BaseDao[T]) DeleteHardById(id string) error {
	return dao.DB.Unscoped().Delete(new(T), "id = ?", id).Error
}

func (dao BaseDao[T]) DeleteByIds(ids ...string) error {
	return dao.DB.Delete(new(T), "id in (?)", ids).Error
}

func (dao BaseDao[T]) DeleteHardByIds(ids ...string) error {
	return dao.DB.Unscoped().Delete(new(T), "id in (?)", ids).Error
}

func (dao BaseDao[T]) DeleteByCondition(conditions *Conditions) error {
	return conditions.To(dao.DB).Delete(new(T)).Error
}

func (dao BaseDao[T]) DeleteHardByCondition(conditions *Conditions) error {
	return conditions.To(dao.DB).Unscoped().Delete(new(T)).Error
}

func (dao BaseDao[T]) Update(record *T) error {
	return dao.DB.Save(record).Error
}

func (dao BaseDao[T]) UpdateSelective(record *T, values any) error {
	return dao.DB.Model(record).Updates(values).Error
}

func (dao BaseDao[T]) UpdateByCondition(record *T, conditions *Conditions) error {
	return conditions.To(dao.DB).Model(new(T)).Updates(structs.ToMapInterface(record)).Error
}

func (dao BaseDao[T]) UpdateSelectiveByCondition(record *T, conditions *Conditions) error {
	return conditions.To(dao.DB).Model(new(T)).Updates(record).Error
}

func (dao BaseDao[T]) SelectById(id string) (*T, error) {
	var record T
	err := dao.DB.Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (dao BaseDao[T]) SelectByIds(ids ...string) ([]*T, error) {
	var records []*T
	err := dao.DB.Where("id = (?)", ids).First(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (dao BaseDao[T]) SelectOne(record *T) (*T, error) {
	var entity T
	err := dao.DB.Where(record).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (dao BaseDao[T]) SelectOneByConditions(record *T, conditions *Conditions) (*T, error) {
	var entity T
	err := conditions.To(dao.DB).Where(record).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (dao BaseDao[T]) Select(record *T) ([]*T, error) {
	var entities []*T
	err := dao.DB.Where(record).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (dao BaseDao[T]) SelectAll() ([]*T, error) {
	var entities []*T
	err := dao.DB.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (dao BaseDao[T]) SelectCount(record *T) (int64, error) {
	var count int64
	err := dao.DB.Model(new(T)).Where(record).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao BaseDao[T]) SelectByCondition(conditions *Conditions) ([]*T, error) {
	var entities []*T
	err := conditions.To(dao.DB).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (dao BaseDao[T]) SelectCountByCondition(conditions *Conditions) (int64, error) {
	var count int64
	err := conditions.To(dao.DB).Model(new(T)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dao BaseDao[T]) SelectDistinct(conditions *Conditions) ([]*T, error) {
	var entities []*T
	err := conditions.To(dao.DB).Distinct().Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (dao BaseDao[T]) SelectPage(currentPage, pageSize int, sort string) (*Pagination[T], error) {
	return dao.SelectPageByCondition(currentPage, pageSize, sort, NewConditions())
}

func (dao BaseDao[T]) SelectPageByCondition(currentPage, pageSize int, sort string, conditions *Conditions) (*Pagination[T], error) {
	db := conditions.To(dao.DB)
	var pagination Pagination[T]
	var entities []*T
	var totalCount int64
	var totalPages int64
	// 计算总记录数
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// 获取总页数
	totalPages = totalCount / int64(pageSize)
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}
	// 当前页
	pageIndex := (currentPage - 1) * pageSize
	err := db.Order(sort).
		Offset(pageIndex).
		Limit(pageSize).
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	pagination.CurrentPage = currentPage
	pagination.TotalCount = totalCount
	pagination.PageSize = pageSize
	pagination.TotalPages = totalPages
	pagination.DataCollection = entities
	return &pagination, nil
}

func (dao BaseDao[T]) Exist(record *T) (bool, error) {
	count, err := dao.SelectCount(record)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
