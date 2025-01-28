package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := &Role{}

	err := s.db.QueryRowContext(ctx, query, slug).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return role, nil
}
