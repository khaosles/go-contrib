package g

import (
	"container/list"

	"gorm.io/gorm"
)

/*
   @File: conditions.go
   @Author: khaosles
   @Time: 2023/6/11 19:53
   @Desc: 条件查询 支持 and or select join group 调用该结构体To方法生成带查询条件的db
*/

type Conditions struct {
	queue *list.List
}

type condition struct {
	query string
	join  int
	value []any
}

const (
	AND = iota
	OR
	ORDER
	SELECT
	GROUP
	JOIN
)

func NewConditions() *Conditions {
	queue := list.New()
	return &Conditions{queue: queue}
}

func (c *Conditions) To(db *gorm.DB) *gorm.DB {
	if c.queue.Len() == 0 {
		return db
	}
	for e := c.queue.Front(); e != nil; e = e.Next() {
		cd := e.Value.(*condition)
		switch cd.join {
		case AND:
			db = db.Where(cd.query, cd.value...)
		case ORDER:
			db = db.Order(cd.query)
		case SELECT:
			db = db.Select(cd.query)
		case GROUP:
			db = db.Group(cd.query)
		case OR:
			db = db.Or(cd.query, cd.value...)
		case JOIN:
			db = db.Joins(cd.query)
		}
	}
	return db
}

// join 连接方式 or and，op 运算符 field 字段  value 字段值
func (c *Conditions) sql(join int, query string, values ...any) *Conditions {
	c.queue.PushBack(&condition{join: join, query: query, value: values})
	return c
}

func (c *Conditions) Order(fields string) *Conditions {
	return c.sql(ORDER, fields)
}

func (c *Conditions) Select(fields string) *Conditions {
	return c.sql(SELECT, fields)
}

func (c *Conditions) Group(fields string) *Conditions {
	return c.sql(GROUP, fields)
}

func (c *Conditions) Joins(joinCondition string) *Conditions {
	return c.sql(JOIN, joinCondition)
}

func (c *Conditions) AndIn(field string, value any) *Conditions {
	return c.sql(AND, field+" in (?)", value)
}

func (c *Conditions) AndNotIn(field string, value any) *Conditions {
	return c.sql(AND, field+" not in (?)", value)
}

func (c *Conditions) AndEqualTo(field string, value any) *Conditions {
	return c.sql(AND, field+" = ?", value)
}

func (c *Conditions) AndLessThan(field string, value any) *Conditions {
	return c.sql(AND, field+" < ?", value)
}

func (c *Conditions) AndLessThanOrEqualTo(field string, value any) *Conditions {
	return c.sql(AND, field+" <= ?", value)
}

func (c *Conditions) AndGreaterThan(field string, value any) *Conditions {
	return c.sql(AND, field+" > ?", value)
}

func (c *Conditions) AndGreaterThanOrEqualTo(field string, value any) *Conditions {
	return c.sql(AND, field+" >= ?", value)
}

func (c *Conditions) AndNotEqualTo(field string, value any) *Conditions {
	return c.sql(AND, field+" <> ?", value)
}

func (c *Conditions) AndLike(field string, value any) *Conditions {
	return c.sql(AND, field+" like ?", value)
}

func (c *Conditions) AndNotLike(field string, value any) *Conditions {
	return c.sql(AND, field+" not like ?", value)
}

func (c *Conditions) AndILike(field string, value any) *Conditions {
	return c.sql(AND, field+" ilike ?", value)
}

func (c *Conditions) AndBetween(field string, value1, value2 any) *Conditions {
	return c.sql(AND, field+" between ? and ?", value1, value2)
}

func (c *Conditions) AndNotBetween(field string, value1, value2 any) *Conditions {
	return c.sql(AND, field+" not between ? and ?", value1, value2)
}

func (c *Conditions) AndIsNull(field string) *Conditions {
	return c.sql(AND, field+" is null")
}

func (c *Conditions) AndIsNotNull(field string) *Conditions {
	return c.sql(AND, field+" is not null")
}

func (c *Conditions) OrIn(field string, value any) *Conditions {
	return c.sql(OR, field+" in (?)", value)
}

func (c *Conditions) OtNotIn(field string, value any) *Conditions {
	return c.sql(OR, field+" not in (?)", value)
}

func (c *Conditions) OrEqualTo(field string, value any) *Conditions {
	return c.sql(OR, field+" = ?", value)
}

func (c *Conditions) OrLessThan(field string, value any) *Conditions {
	return c.sql(OR, field+" < ?", value)
}

func (c *Conditions) OrLessThanOrEqualTo(field string, value any) *Conditions {
	return c.sql(OR, field+" <= ?", value)
}

func (c *Conditions) OrGreaterThan(field string, value any) *Conditions {
	return c.sql(OR, field+" > ?", value)
}

func (c *Conditions) OrGreaterThanOrEqualTo(field string, value any) *Conditions {
	return c.sql(OR, field+" >= ?", value)
}

func (c *Conditions) OrNotEqualTo(field string, value any) *Conditions {
	return c.sql(OR, field+" <> ?", value)
}

func (c *Conditions) OrLike(field string, value any) *Conditions {
	return c.sql(OR, field+" like ?", value)
}

func (c *Conditions) OrNotLike(field string, value any) *Conditions {
	return c.sql(OR, field+" not like ?", value)
}

func (c *Conditions) OrILike(field string, value any) *Conditions {
	return c.sql(OR, field+" ilike ?", value)
}

func (c *Conditions) OrIsNull(field string) *Conditions {
	return c.sql(OR, field+" is null")
}

func (c *Conditions) OrIsNotNull(field string) *Conditions {
	return c.sql(OR, field+" is not null")
}

func (c *Conditions) OrBetween(field string, value1, value2 any) *Conditions {
	return c.sql(OR, field+" between ? and ?", value1, value2)
}

func (c *Conditions) OrNotBetween(field string, value1, value2 any) *Conditions {
	return c.sql(OR, field+" not between ? and ?", value1, value2)
}
