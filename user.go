package site

import "github.com/jackc/pgx/v4/pgxpool"

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{pool: pool}
}

// func (svc *UserSvc)authenticateUser(email string, password string) int, bool {

// }
