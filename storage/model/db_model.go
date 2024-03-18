package model

import (
	"time"
)

type ConfigParam struct {
	ID       uint      `gorm:"primaryKey;autoIncrement;notnull"`
	Name     string    `gorm:"type:varchar(20);not null;unique;"`
	Value    string    `gorm:"type:varchar(255);not null;"`
	UpdateAt time.Time `gorm:"autoCreateTime;"`
}

// OpusNode 分形章鱼的节点
type OpusNode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;notnull"`
	Uri       string    `gorm:"index;type:varchar(2048);not null;unique;"`
	Owner     string    `gorm:"index;type:varchar(128);"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

// OpusStream 每个节点的上游及其分润比例
type OpusStream struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;notnull"`
	CurrUriId   uint      `gorm:"index;not null"`
	UpstreamUri string    `gorm:"index;type:varchar(2048);not null"`
	Ratio       float64   `gorm:"index;type:decimal(6,4);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;"`
}
