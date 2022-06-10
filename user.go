package site

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (svc *UserService) CreateAdmin(ctx context.Context, email string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Could not hash password. %w", err)
	}
	_, err = svc.db.Exec(ctx, "INSERT INTO users (email, password, admin) VALUES ($1, $2, true)", email, hash)
	if err != nil {
		return fmt.Errorf("Could not insert record. %w", err)
	}
	return nil
}

// func (svc *UserSvc)authenticateUser(email string, password string) int, bool {

// }
