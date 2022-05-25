package swiki

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func (svc *UserService) createTestUser(ctx context.Context, t *testing.T) int {
	t.Helper()
	var userid int
	email := "john@doe.com"
	password, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Could not create test user. %s", err)
	}
	err = svc.pool.QueryRow(ctx, "insert into users (email, password) values ($1, $2) returning user_id", email, password).Scan(&userid)
	if err != nil {
		t.Fatalf("Could not create test user. %s", err)
	}
	return userid
}
