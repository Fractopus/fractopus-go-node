package model

import (
	"time"
)

type OpusUri struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;notnull"`
	Uri       string    `gorm:"index;type:varchar(2048);not null;unique;"`
	CreatedAt time.Time `gorm:"autoCreateTime;"`
}

type OpusStream struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;notnull"`
	CurrUriId     uint      `gorm:"index;not null"`
	UpstreamUriId uint      `gorm:"index;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime;"`
}
