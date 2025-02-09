package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken table model..
type RefreshToken struct {
	ID           uint      `gorm:"primaryKey"`
	OwnedBy      uuid.UUID `gorm:"type:char(36);not null"`
	RefreshToken string    `gorm:"type:varchar(512);not null"`
	HostInfo     string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	User         Account `gorm:"foreignKey:OwnedBy"`
}

type RefreshTokenInput struct {
	OwnedBy      uuid.UUID `gorm:"type:char(36);not null"`
	RefreshToken string    `gorm:"type:varchar(512);not null"`
	HostInfo     string    `gorm:"type:varchar(255);not null"`
}

func SaveRefreshToken(uuid uuid.UUID, refreshTokenString string, hostInfo string) error {
	input := RefreshToken{}
	input.OwnedBy = uuid
	input.RefreshToken = refreshTokenString
	input.HostInfo = hostInfo
	_, err := input.SaveRefreshTokenDB()
	if err != nil {
		return err
	}
	return nil
}

func (a *RefreshToken) SaveRefreshTokenDB() (*RefreshToken, error) {
	err := DB.Create(&a).Error
	if err != nil {
		return &RefreshToken{}, err
	}
	return a, nil
}

func RemoveRefreshTokenDB(hostInfo string) error {
	token := RefreshToken{}
	DB.Where("host_Info = ?", hostInfo).Delete(&token)
	return nil
}
