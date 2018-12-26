package models

type Merchant struct {
	ID         int
	MerchantID string
	Status     string
}

const MERCHANT_STATUS_ACTIVE = "ACTIVE"
const MERCHANT_STATUS_INACTIVE = "INACTIVE"

func (merchant *Merchant) TableName() string {
	return "merchants"
}
