package swiki

import (
	"context"
	"testing"
)

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
