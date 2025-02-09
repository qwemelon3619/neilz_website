package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	UserUUID     uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID       string    `gorm:"type:varchar(255);unique;not null"`
	UserPassword string    `gorm:"type:varchar(512);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}
type AccountInput struct {
	UserID       string
	UserPassword string
}

func SaveUserEmailAndPassword(id string, password string) error {
	var account = Account{}
	account.UserUUID = uuid.New()
	account.UserID = id
	account.UserPassword = password

	_, err := account.SaveAccountDB()
	if err != nil {
		return errors.New("DataBase Error: could not save the data")
	}
	return nil
}

func (a *Account) SaveAccountDB() (*Account, error) {
	err := DB.Create(&a).Error
	if err != nil {
		return &Account{}, err
	}
	return a, nil
}

func FindUserByUserID(id string) (Account, error) {
	var account Account
	result := DB.First(&account, "user_id = ?", id)
	if result.Error != nil {
		return Account{}, errors.New("DB error")
	}
	return account, nil
}
