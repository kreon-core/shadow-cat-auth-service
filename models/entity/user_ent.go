package entity

import "github.com/google/uuid"

type User struct {
	TimeMeta

	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(128);uniqueIndex"`
	PasswordHash string    `gorm:"type:text"`

	Role        string `gorm:"type:varchar(32);not null"`
	DisplayName string `gorm:"type:varchar(128)"`
	AvatarURL   string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(32);not null;default:'active'"`

	PlayerID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
}
