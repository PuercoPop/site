package swiki

import (
	"crypto/rand"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Authenticate checks the user credentials. If valid the user id is returned. If not an error is returned.
func Authenticate(email string, password string) (int, error) {
	return 0, nil
	// var hash []byte
	// row := db.QueryRow("SELECT password FROM users where email = $1", email)
	// err = row.Scan(&hash) // TODO:: Check
}

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("Could not read random bytes. %s", err)
	}
	return buf

}

type SessionService struct {
	pool *pgxpool.Pool
}

func NewSessionService(pool *pgxpool.Pool) *SessionService {
	return &SessionService{pool: pool}
}

// CreateSessionFor creates
func (svc *SessionService) CreateSessionFor(user_id int) ([]byte, error) {
	// "The session ID should be at least 128 bits o prevent brute-force session guessing attacks."
	// Ref: https://owasp.org/www-community/vulnerabilities/Insufficient_Session-ID_Length
	sessionid := randomBytes(128)
	// insert to sql
	return sessionid, nil
}

// UserFromSession retrieves the user-id associated with session id, sid.
func (svc *SessionService) UserFromSession(sid []byte) (int, error) {
	return 0, nil
}
