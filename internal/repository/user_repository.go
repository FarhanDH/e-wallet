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

	// Debit Sender
	res, err := tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromID)
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

	// Credit receiver
	res, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toID)
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

	// Record history
	_, err = tx.Exec("INSERT INTO transactions (from_account_id, to_account_id, amount, transaction_type) VALUES ($1, $2, $3, 'TRANSFER')", fromID, toID, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetTransactionHistory(userID int) ([]entity.Transaction, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, transaction_type, created_at 
	FROM transactions
	WHERE from_account_id = $1 OR to_account_id = $1
	ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []entity.Transaction

	for rows.Next() {
		var t entity.Transaction
		err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.TransactionType, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, t)
	}
	return history, nil
}
