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

type OpusUri struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;notnull"`
	Uri       string    `gorm:"index;type:varchar(2048);not null;unique;"`
	Owner     string    `gorm:"index;type:varchar(128);not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

type OpusStream struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;notnull"`
	CurrUriId     uint      `gorm:"index;not null"`
	UpstreamUriId uint      `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime;"`
}
