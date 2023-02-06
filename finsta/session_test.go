package finsta

// import (
// 	"context"
// 	"net/http"
// 	"testing"
// )

// // TODO(javier): Test SessionMiddleware
// type SessionMemStore struct {
// 	// a map from []byte to it
// 	m map[string]int
// }

// func NewSessionMemStore() *SessionMemStore {
// 	store := &SessionMemStore{}
// 	store.m = make(map[string]int)
// 	return store

// }
// func (store *SessionMemStore) New(ctx context.Context, email string, password string) ([]byte, error) {
// 	sid := randomBytes(128)
// 	// How do I get th euser id?
// 	return sid, nil
// }
// func (store *SessionMemStore) Lookup(r *http.Request) (int, error) {
// 	return 0, nil
// }
// func TestSessionMiddleware(t *testing.T) {

// }

// func sqltest.DB(ctx context, t *testing.T) (*pgx.Conn, func()) {
// 	t.Helper()
// 	_, ok := os.LookupEnv("DBURL")
// 	if !ok {
// 		t.Skipf("%t requries")
// 	}
// 	conf, err := pgxpool.ParseConfig(url)
// 	pool, db := pgxpool.ConnectConfig(ctx, config)
// 	close := func() {
// 		conn.Close()
// 	}
// 	return db, close

// }

// func TestRetrieveUserFromSession(t *testing.T) {
// 	ctx := context.Background()
// 	db, close := sqltest.DB(t)
// 	defer close()
// 	svc := NewSessionService(db)
// 	usersvc := NewUserService(db)
// 	want := usersvc.createTestUser(ctx, t)
// 	sid, err := svc.CreateSessionFor(ctx, want)
// 	if err != nil {
// 		t.Fatalf("Could not create session. %s", err)
// 	}
// 	// TODO(javier): check the expiration time on the recently created session
// 	got, err := svc.UserFromSession(ctx, sid)
// 	if err != nil {
// 		t.FailNow()
// 	}
// 	if got != want {
// 		t.Errorf("User ids don't match. got: %d. want: %d.", got, want)
// 	}

// }
