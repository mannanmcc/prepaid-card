package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//NewDB - NewDB returns an instance of gorm database connection
func NewDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}

	return db, nil
}
