package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func RunInTransaction(conn *gorm.DB, fn func(*gorm.DB) error) error {
	tx := conn.Begin()
	if err := tx.Error; err != nil {
		return errors.New(err.Error())
	}
	needRollback := true
	defer func() {
		if needRollback && tx != nil {
			tx.Rollback()
		}
	}()
	if err := fn(tx); err != nil {
		return errors.New(err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return errors.New(err.Error())
	}
	needRollback = false
	return nil
}
