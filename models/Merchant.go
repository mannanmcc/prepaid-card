package models

type Merchant struct {
	ID         int
	MerchantID string
	Status     string
}

const MerchantStatusActive = "ACTIVE"
const MerchantStatusInactive = "INACTIVE"

func (merchant *Merchant) TableName() string {
	return "merchants"
}
