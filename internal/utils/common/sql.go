package common

import "database/sql"

type TxFunc func(tx *sql.Tx) error

// ExecuteInTransaction - выполняет функцию в рамках транзакции
func ExecuteInTransaction(db *sql.DB, fn TxFunc) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return nil
}
