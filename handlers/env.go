package handlers

import "github.com/jinzhu/gorm"

type Env struct {
	Db *gorm.DB
}

const ACCOUNT_STATUS_ACTIVE = "ACTIVE"
