package repository

import (
	"database/sql"
	"ewallet/internal/entity"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(name string, initialBalance int) error {
	query := "INSERT INTO users (name, balance) VALUES ($1, $2)"

	_, err := r.db.Exec(query, name, initialBalance)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUser(id int) (entity.User, error) {
	var user entity.User

	query := "SELECT id, name, balance FROM users WHERE id = $1"

	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Balance)

	if err != nil && err == sql.ErrNoRows {
		return entity.User{}, fmt.Errorf("Not Found")
	}

	return user, err
}

func (r *UserRepository) Transfer(fromID int, toID int, amount int) error {
	// Begin Transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Reduce the balance of sender
	queryDebit := "UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1"
	res, err := tx.Exec(queryDebit, amount, fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Rollback when there's no rows affected, in this case, the balance is not enough or the user is not found
	rows, _ := res.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("insufficient funds or sender not found")
	}

	// Add the balance of receiver
	queryCredit := "UPDATE users SET balance = balance + $1 WHERE id = $2"
	res, err = tx.Exec(queryCredit, amount, toID)

	if err != nil {
		tx.Rollback()
		return err
	}
	// Rollback when the receiver is not found
	rows, _ = res.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("receiver with ID %d not found", toID)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
