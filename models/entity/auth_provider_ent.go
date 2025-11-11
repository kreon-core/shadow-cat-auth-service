package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuthProvider struct {
	TimeMeta

	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Provider    string         `gorm:"type:varchar(64);not null"`
	ProviderUID string         `gorm:"type:varchar(256);not null"`
	LinkedAt    time.Time      `gorm:"type:timestamptz;not null"`
	Metadata    datatypes.JSON `gorm:"type:jsonb"`

	// Associations
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
