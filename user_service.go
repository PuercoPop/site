package swiki

import "github.com/jackc/pgx/v4/pgxpool"

type UserSvc struct {
	db *pgxpool.Pool
}

func (svc *UserSvc)authenticateUser(email string, password string) int, bool {

}
