package pgx

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.pool.QueryRow(context.Background(),
		`SELECT id, email, password_hash, status, created_at, updated_at
		 FROM users WHERE email = $1`, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	row := r.pool.QueryRow(context.Background(),
		`SELECT id, email, password_hash, status, created_at, updated_at
		 FROM users WHERE id = $1`, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Email, user.PasswordHash, user.Status, user.CreatedAt, user.UpdatedAt)
	return err
}
