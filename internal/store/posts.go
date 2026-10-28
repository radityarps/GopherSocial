package store

import (
	"context"
	"database/sql"
)

type Post struct {
	ID        int    `json: "id"`
	Content   string `json: "content"`
	UserID    int    `json: "user_id"`
	CreatedAt string `json: "created_at"`
	UpdatedAt string `json: "updated_at"`
}

type PostgresPostStore struct {
	db *sql.DB
}

func (s *PostgresPostStore) Create(ctx context.Context, post *Post) error {
	err := s.db.QueryRowContext(ctx, "INSERT INTO posts (content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, updated_at, created_at", post.Content, post.UserID, post.CreatedAt, post.UpdatedAt).Scan(&post.ID, &post.UpdatedAt, &post.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
