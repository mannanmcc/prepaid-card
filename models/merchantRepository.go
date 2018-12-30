package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//MerchantRepositoryInterface - repository interface for merchant
type MerchantRepositoryInterface interface {
	FindByMerhcantID(merchantID string) (*Merchant, error)
}

//MerchantRepository - type for account repository
type MerchantRepository struct {
	Db *gorm.DB
}

//FindByMerchantID - returns merchant
func (repo *MerchantRepository) FindByMerchantID(merchantID string) (Merchant, error) {
	var merchant Merchant
	res := repo.Db.Find(&merchant, &Merchant{MerchantID: merchantID})

	if res.RecordNotFound() {
		return merchant, fmt.Errorf("merchant not found with merchant id : %s", merchantID)
	}

	return merchant, nil
}
