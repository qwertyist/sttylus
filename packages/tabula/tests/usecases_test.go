package usecases_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/qwertyist/tabula/repo"
	"github.com/qwertyist/tabula/session"
)

var s session.SessionService
var original session.Session
var now time.Time
var t1 *session.User
var t2 *session.User
var c1 *session.User
var c2 *session.User
var c3 *session.User
var id1 string
var id2 string
var id3 string

func init() {
	db := repo.NewRepository(repo.OpenBoltDB("test.db"))
	s = session.NewSessionService(db)
	t1 = s.CreateUser(session.User{
		Name:    "David Johansson",
		Email:   "david@skrivasnabbt.se",
		Creator: true,
	})

	t2 = s.CreateUser(session.User{
		Name:  "Annan Usersson",
		Email: "",
		Phone: "090-770000",
	})

	c1 = s.CreateUser(session.User{
		Name:  "Anonym tolkanvändare",
		Email: "",
	})
	c2 = s.CreateUser(session.User{
		Name:  "Anonym tolkanvändare",
		Email: "",
	})
	c3 = s.CreateUser(session.User{
		Name:  "Anonym tolkanvändare",
		Email: "",
	})

	now = time.Now()
	original = session.Session{
		ID:          "12341234",
		Name:        "Testsession",
		Recurring:   false,
		From:        now,
		To:          now.Add(time.Hour),
		Type:        session.Tabula,
		SaveText:    true,
		Ref:         t2.ID,
		Itprs:       []string{t1.ID},
		Clients:     []string{c1.ID},
		Description: "Testing testing",
	}
	for i := 0; i < 7; i++ {
		sess := original
		if i%2 == 0 {
			sess.Clients = append(sess.Clients, c2.ID)
		}
		if i%3 == 0 {
			sess.Clients = append(sess.Clients, c3.ID)
		}
		sess.ID = ""
		created, _ := s.CreateSession(sess)
		if id1 == "" {
			id1 = created.ID
		} else if id2 == "" {
			id2 = created.ID
		} else if id3 == "" {
			id3 = created.ID
		}

	}
}
func TestCreateSession(t *testing.T) {
	test := original
	t.Run("from_to_opposite", func(t *testing.T) {
		test.From = now.Add(time.Hour)
		test.To = now
		created, err := s.CreateSession(test)
		if err != nil {
			if err.Error() != session.InvalidTimeToBeforeFrom {
				t.Errorf("got %q, wanted %q", err, session.InvalidTimeToBeforeFrom)
			}
			return
		}
		if created != nil {
			t.Errorf("session created, should've failed due to %q:", session.InvalidTimeToBeforeFrom)
		}

	})

	t.Run("from_before_today", func(t *testing.T) {
		test := original
		test.To = now.AddDate(0, 0, -2)
		test.From = now.Add(time.Duration(-5)*time.Hour).AddDate(0, 0, -2)
		created, err := s.CreateSession(test)
		want := session.InvalidTimeFromBeforeToday
		if err != nil {
			got := err.Error()
			if got != want {
				t.Errorf("got %q, wanted %q", got, want)
			}
		} else {
			if created != nil {
				t.Errorf("session created, should've failed due to %q:", want)
			}
		}

	})

	t.Run("should_create_placeholder_name", func(t *testing.T) {
		test := original
		test.ID = ""
		test.Name = ""
		test.From = time.Now()
		test.To = time.Now().Add(1 * time.Hour)
		created, err := s.CreateSession(test)
		if err != nil {
			t.Errorf("error %q:", err.Error())
		}
		if has := strings.HasSuffix(created.Name, created.ID); !has {
			t.Errorf("session placeholder doesn't follow pattern '{type} {ID}'")
		}
	})

	t.Run("should_have_atleast_one_person", func(t *testing.T) {
		test := original
		test.Clients = []string{}
		test.Ref = ""
		test.Itprs = nil

		created, err := s.CreateSession(test)
		want := session.WarningNoReference
		if created == nil {
			t.Errorf("creation failed, for no known reason")
		}
		if err != nil {
			if err.Error() != want {
				t.Errorf("error %q:", err.Error())
			}
		}
	})

	t.Run("should_not_create_duplicate_id", func(t *testing.T) {

	})

}

func TestGetSession(t *testing.T) {
	t.Run("should_return_nothing_on_empty_ID", func(t *testing.T) {
		found, _ := s.GetSession("")
		if found != nil {
			t.Errorf("shouldn't return anything")
		}
	})

}

func TestGetSessions(t *testing.T) {
	t.Run("no_ID_should_return_all_lists", func(t *testing.T) {
		found, err := s.GetSessions([]string{""})
		if err != nil {
			t.Errorf("should return list of sessions")
		}
		if found == nil {
			t.Errorf("should return list of sessions")
		}
		found, _ = s.GetSessions(nil)
		if found == nil {
			t.Errorf("should return list of sessions")
		}
	})

	t.Run("providing_person_id_returns_all_linked_sessions", func(t *testing.T) {
		found, _ := s.GetSessions([]string{c2.ID})
		if len(found) != 4 {
			t.Errorf("should return 4 sessions")
		}
	})
}

func TestUpdateSession(t *testing.T) {
	rand.Seed(now.Unix())
	ss, err := s.GetSessions(nil)
	if err != nil {
		t.Errorf("something when wrong when retrieving all sessions: %q", err)
	}

	var uu []*session.Session
	is := make(map[int]bool)

	for i := 0; i < 5; i++ {
		n := rand.Int() % len(ss)
		if !is[n] {
			is[n] = true
			uu = append(uu, ss[n])
		}
	}

	t.Run("prevent_changing_session_id", func(t *testing.T) {
		u := uu[0]
		uu = append(uu[:0], uu[1:]...)
		u.ID = "12345678"
		updated, _ := s.UpdateSession(*u)
		if updated != nil {
			t.Errorf("shouldn't be able to update a session ID")
		}
	})
	t.Run("can_bulk_update_with_random_values", func(t *testing.T) {
		for _, u := range uu {
			updated, err := s.UpdateSession(*u)
			if err != nil {
				t.Errorf("failed bulk updating due to: %q", err.Error())
			}
			if updated == nil {
				t.Errorf("couldn't update session: %+v", u)
			}
		}

	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("DeleteASession", func(t *testing.T) {
		sess, err := s.DeleteSession(original.ID)
		if sess == nil {
			t.Errorf("couldn't delete session")
		}
		if err != nil {
			t.Errorf("couldn't delete session due to: %q", err)
		}
		sess, err = s.DeleteSession(original.ID)
		if sess != nil {
			t.Errorf("shoulnd't be able to delete the same session again")
		}
		if err == nil {
			t.Errorf("shoulnd't be able to delete the same session again")
		}

	})
}
