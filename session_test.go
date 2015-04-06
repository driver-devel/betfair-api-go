package betfair

import (
	"log"
	"testing"
)

var testSession *Session

func GetTestSession() *Session {
	if testSession == nil {
		session, err := NewSession(GetTestAccount())

		if err != nil {
			log.Panic(err)
		}

		testSession = session
	}

	return testSession
}

func TestSessionInteractiveLogin(t *testing.T) {
	session := GetTestSession()
	token, err := session.GetToken()

	if err != nil {
		t.Error(err)
		return
	}

	if token == "" {
		t.Error("Token is blank")
	}
}

func TestKeepAlive(t *testing.T) {
	session := GetTestSession()

	ok, err := session.KeepAlive()

	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("Keep alive failure")
	}
}
