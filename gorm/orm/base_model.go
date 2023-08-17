package g

import (
	"github.com/khaosles/giz/gen"
	"gorm.io/gorm"
)

/*
   @File: base_model.go
   @Author: khaosles
   @Time: 2023/6/11 11:10
   @Desc: model 继承该类，则拥有常见基础属性，以及自动生成无分隔符的uuid
*/

type AuthMixin struct {
	CreateBy string `json:"createBy,omitempty" gorm:"column:create_by;type:varchar(100);default:null;comment:创建人"` // 创建人
	UpdateBy string `json:"updateBy,omitempty" gorm:"column:update_by;type:varchar(100);default:null;comment:更新人"` // 更新人
}

type PgTimeMixin struct {
	CreateTime Timestamp      `json:"createTime,omitempty" gorm:"autoCreateTime;column:create_time;type:timestamptz;comment:创建时间"` // 创建时间
	UpdateTime Timestamp      `json:"updateTime,omitempty" gorm:"autoUpdateTime;column:update_time;type:timestamptz;comment:更新时间"` // 更新时间
	DeleteTime gorm.DeletedAt `json:"-" gorm:"index;column:delete_time;type:timestamptz;comment:删除时间"`                             // 删除标记
}

type TimeMixin struct {
	CreateTime Timestamp      `json:"createTime,omitempty" gorm:"autoCreateTime;column:create_time;type:datetime;comment:创建时间"` // 创建时间
	UpdateTime Timestamp      `json:"updateTime,omitempty" gorm:"autoUpdateTime;column:update_time;type:datetime;comment:更新时间"` // 更新时间
	DeleteTime gorm.DeletedAt `json:"-" gorm:"index;column:delete_time;type:datetime;comment:删除时间"`                             // 删除标记
}

type IdMixin struct {
	ID string `json:"id" gorm:"primaryKey;column:id;type:varchar(32);comment:主键"` // 主键ID
}

func (m *IdMixin) BeforeCreate(tx *gorm.DB) error {
	//node, _ := snowflake.NewNode(0)
	//m.ID = node.Generate().String()
	m.ID = gen.UuidNoSeparator()
	return nil
}

type BaseModel struct {
	IdMixin
	PgTimeMixin
	AuthMixin
	Remarks string `json:"remarks,omitempty" gorm:"column:remarks;default:null;comment:备注"` // 备注
}

type BaseModelNoPg struct {
	IdMixin
	TimeMixin
	AuthMixin
	Remarks string `json:"remarks,omitempty" gorm:"column:remarks;default:null;comment:备注"` // 备注
}
