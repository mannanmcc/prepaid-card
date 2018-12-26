package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

//TopupRepositoryInterface - repository interface for Topup
type TopupRepositoryInterface interface {
	Store(topup Topup) (*Topup, error)
}

//TopupRepository - type for topup repository
type TopupRepository struct {
	Db *gorm.DB
}

//Store - save the topup entity
func (repo *TopupRepository) Store(topup *Topup) error {
	id := repo.Db.Save(&topup)

	if id == nil {
		return errors.New("topup saving failed")
	}

	return nil
}

//FindAllTopups - return all topups by card number
func (repo *TopupRepository) FindAllTopups(cardNumber string) []Topup {
	var topups []Topup
	res := repo.Db.Find(&topups, &Topup{CardNumber: cardNumber})

	if res.RecordNotFound() {
		return nil
	}

	return topups
}
