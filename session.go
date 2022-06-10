package site

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("Could not read random bytes. %s", err)
	}
	return buf

}

type SessionHandler struct {
	db *pgxpool.Pool
}

func NewSessionService(pool *pgxpool.Pool) *SessionHandler {
	return &SessionHandler{db: pool}
}

// Authenticate checks the user credentials. If valid the user id is returned. If not an error is returned.
func (svc *SessionHandler) Authenticate(ctx context.Context, email string, password string) ([]byte, error) {
	var userid int
	var hashed []byte
	err := svc.db.QueryRow(ctx, "SELECT user_id, password FROM users where email = $1", email).Scan(&userid, &hashed)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return nil, err
	}
	// "The session ID should be at least 128 bits o prevent brute-force session guessing attacks."
	// Ref: https://owasp.org/www-community/vulnerabilities/Insufficient_Session-ID_Length
	sid := randomBytes(128)
	_, err = svc.db.Exec(ctx, "INSERT INTO sessions (session_id, user_id) values ($1, $2)", sid, userid)
	if err != nil {
		return nil, err
	}
	return sid, nil
}

// SessionMiddleware reads the session id cookie and adds the associated user to
// the request context.
func SessionMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
