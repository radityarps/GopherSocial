package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) Create(ctx context.Context, user *User) error {
	err := s.db.QueryRowContext(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at", user.Name, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
