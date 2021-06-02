package database

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("not found record")
	ErrKeyConflict    = errors.New("conflict key")
	ErrFKConstraint   = errors.New("a foreign key constraint fails")
)

// WrapError wrap database error to handle cause.
// ErrRecordNotFound is returned if err is a gorm.ErrRecordNotFound
// If conflict key error, then ErrKeyConflict will return.
func WrapError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return ErrRecordNotFound
	}
	if e, ok := err.(*mysql.MySQLError); ok {
		switch e.Number {
		case 1062:
			return ErrKeyConflict
		case 1452:
			return ErrFKConstraint
		}
	}
	return err
}
