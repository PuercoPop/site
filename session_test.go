package swiki

import "testing"

func TestRetrieveUserFromSession(t *testing.T) {
	db, close := setuptestdb(t)
	defer close()
	svc := NewSessionService(db)
	want := 19
	sid, err := svc.CreateSessionFor(want)
	if err != nil {
		t.FailNow()
	}
	got, err := svc.UserFromSession(sid)
	if err != nil {
		t.FailNow()
	}
	if got != want {
		t.Errorf("User ids don't match. got: %d. want: %d.", got, want)
	}
}
