package models

type Merchant struct {
	ID         int `gorm:"primary_key"`
	MerchantID string
	Status     string
}

const MerchantStatusActive = "ACTIVE"
const MerchantStatusInactive = "INACTIVE"

func (merchant *Merchant) TableName() string {
	return "merchants"
}
