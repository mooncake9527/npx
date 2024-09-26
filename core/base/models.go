package base

import (
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy uint64 `json:"createBy" gorm:"type:bigint unsigned;index;comment:创建者"` //创建者id
	UpdateBy uint64 `json:"updateBy" gorm:"type:bigint unsigned;index;comment:更新者"` //更新者id
}

type StatusModel struct {
	Status int `json:"status" gorm:"type:tinyint unsigned;comment:状态 1 默认状态 2 成功 3 失败"`
}

type Model struct {
	Id int `json:"id" gorm:"type:bigint unsigned;primaryKey;autoIncrement;comment:主键编码"` //主键
}

type ModelIntTime struct {
	CreatedAt time.Time  `json:"createdAt" gorm:"type:datetime;comment:创建时间"`   //创建时间戳
	UpdatedAt *time.Time `json:"updatedAt" gorm:"type:datetime;comment:最后更新时间"` //更新时间戳
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`   //创建时间
	UpdatedAt *time.Time     `json:"updatedAt" gorm:"comment:最后更新时间"` //更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`     //删除时间
}
