// Package session implements authentication and session handling
package session

// Authenticate checks the user credentials. If valid the user id is returned. If not an error is returned.
func Authenticate(email string, password string) (int, error) {
	var hash []byte
	row := db.QueryRow("SELECT password FROM users where email = $1", email)
	err = row.Scan(&hash) // TODO:: Check
}
