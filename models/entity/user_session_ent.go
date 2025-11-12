package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSession struct {
	TimeMeta

	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	Token     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	Revoked   bool      `gorm:"type:boolean;not null;default:false"`

	DeviceInfo datatypes.JSON `gorm:"type:jsonb"`
	IPAddress  string         `gorm:"type:varchar(45)"`

	// Associations
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
