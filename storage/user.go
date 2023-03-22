package storage

import (
	"context"
	"fmt"
)

type User struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Balance   float64 `json:"balance"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (s *Storage) CreateUser(ctx context.Context, user CreateUserRequest) (*User, error) {
	rows := s.conn.QueryRowContext(ctx, "INSERT INTO account(firstName,lastName) VALUES($1, $2) RETURNING id, firstName, lastName", user.FirstName, user.LastName)
	return ScanUser(rows)
}

func (s *Storage) GetUser(ctx context.Context, user User) (*User, error) {
	rows, err := s.conn.QueryContext(ctx, "SELECT id, firstName, lastName, balance FROM account WHERE firstName=$1 AND lastName=$2 ", user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	u, err := ScanUser(rows)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, id string, balance float64) (*User, error) {
	rows, err := s.conn.QueryContext(ctx, "UPDATE account SET balance=$1 RETURNING id, firstname, lastname, balance", balance)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	u, err := ScanUser(rows)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*User, error) {
	rows, err := s.conn.QueryContext(ctx, "SELECT id, firstName, lastName, balance FROM account WHERE id=$1 ", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	u, err := ScanUser(rows)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ScanUser(s Scanner) (*User, error) {
	user := &User{}
	if err := s.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Balance); err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return user, nil
}
