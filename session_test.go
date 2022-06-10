package site

import (
	"context"
	"net/http"
	"testing"
)

// TODO(javier): Test SessionMiddleware
type SessionMemStore struct {
	// a map from []byte to it
	m map[string]int
}

func NewSessionMemStore() *SessionMemStore {
	store := &SessionMemStore{}
	store.m = make(map[string]int)
	return store

}
func (store *SessionMemStore) New(ctx context.Context, email string, password string) ([]byte, error) {
	sid := randomBytes(128)
	// How do I get th euser id?
	return sid, nil
}
func (store *SessionMemStore) Lookup(r *http.Request) (int, error) {}
func TestSessionMiddleware(t *testing.T) {

}

func TestRetrieveUserFromSession(t *testing.T) {
	ctx := context.Background()
	db, close := setuptestdb(t)
	defer close()
	svc := NewSessionService(db)
	usersvc := NewUserService(db)
	want := usersvc.createTestUser(ctx, t)
	sid, err := svc.CreateSessionFor(ctx, want)
	if err != nil {
		t.Fatalf("Could not create session. %s", err)
	}
	// TODO(javier): check the expiration time on the recently created session
	got, err := svc.UserFromSession(ctx, sid)
	if err != nil {
		t.FailNow()
	}
	if got != want {
		t.Errorf("User ids don't match. got: %d. want: %d.", got, want)
	}

}
