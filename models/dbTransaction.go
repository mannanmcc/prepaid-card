package models

import "github.com/jinzhu/gorm"

type InTransaction func(tx *gorm.DB) error

func DoInTransaction(fn InTransaction, fn2 InTransaction, db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := fn(tx)
	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			return xerr
		}
		return err
	}

	err = fn2(tx)
	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			return xerr
		}
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
