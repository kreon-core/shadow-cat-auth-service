package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	TimeMeta

	ID           uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string               `gorm:"type:varchar(64);uniqueIndex;not null"`
	Email        datatypes.NullString `gorm:"type:varchar(128);uniqueIndex"`
	PasswordHash datatypes.NullString `gorm:"type:text"`

	Role        string               `gorm:"type:varchar(32);not null"`
	DisplayName datatypes.NullString `gorm:"type:varchar(128)"`
	AvatarURL   datatypes.NullString `gorm:"type:text"`
	Status      string               `gorm:"type:varchar(32);not null;default:'active'"`
}
