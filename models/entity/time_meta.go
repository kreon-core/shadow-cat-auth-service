package entity

import (
	"time"

	"gorm.io/gorm"
)

type TimeMeta struct {
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
