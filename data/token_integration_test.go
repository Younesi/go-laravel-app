//go:build integration

// run tests with this command: go test --cover . -v --tags integration --count=1
package data_test

import (
	"myapp/data"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var dummyUser = data.User{
	FirstName: "Test",
	LastName:  "Guy",
	Email:     "me@dummy.net",
	IsActive:  1,
	Password:  "password",
}

var dummyToken = data.Token{
	FirstName: "Some",
	Email:     "me@here.com",
}

func TestToken_Table(t *testing.T) {
	s := models.Tokens.Table()
	if s != "tokens" {
		t.Error("wrong table name")
	}
}

func TestToken_GenerateToken(t *testing.T) {
	id, err := models.Users.Insert(dummyUser)
	if err != nil {
		t.Error("error inserting user: ", err)
	}

	_, err = models.Tokens.GenerateToken(id, time.Minute*2)
	if err != nil {
		t.Error("err generating token: ", err)
	}
}

func TestToken_Insert(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}
	token, err := models.Tokens.GenerateToken(u.ID, time.Minute*2)
	if err != nil {
		t.Error("failed to generate token: ", err)
	}

	err = models.Tokens.Insert(*token, *u)
	if err != nil {
		t.Error("failed to insert the token")
	}
}

func TestToken_GetUserForToken(t *testing.T) {
	token := "invalid"
	_, err := models.Tokens.GetUserForToken(token)
	if err == nil {
		t.Error("expected error but not recieved when getting user with a bad token")
	}

	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetUserForToken(u.Token.PlainText)
	if err != nil {
		t.Error("failed to get user with valid token: ", err)
	}
}

func TestToken_GetTokensForUser(t *testing.T) {
	u, _ := models.Users.GetByEmail(dummyUser.Email)

	tokens, err := models.Tokens.GetTokensForUser(u.ID)
	if err != nil {
		t.Error("failed to get token for a user: ", err)
	}

	if len(tokens) == 0 {
		t.Error("expected to get tokens of a user, but got none")
	}
}

func TestToken_Get(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.Get(u.Token.ID)
	if err != nil {
		t.Error("error getting token: ", err)
	}
}

func TestToken_GetByToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetByToken(u.Token.PlainText)
	if err != nil {
		t.Error("error getting token: ", err)
	}

	_, err = models.Tokens.GetByToken("invalid token")
	if err == nil {
		t.Error("expected to get error but got none when getting by invalid token")
	}
}

var authData = []struct {
	name          string
	token         string
	email         string
	errorExpected bool
	message       string
}{
	{"invalid", "abcdefghijklmnopqrstuvwxyz", "a@there.com", true, "invalid token accepted as valid"},
	{"invalid_length", "abcdef", "a@there.com", true, "token of wrong length token accepted as valid"},
	{"no_user", "abcdefghijklmnopqrstuvwxyz", "a@there.com", true, "no user, but token accepted as valid"},
	{"valid", "", dummyUser.Email, false, "valid token reported as invalid"},
}

func TestToken_AuthenticateToken(t *testing.T) {
	for _, tt := range authData {
		token := ""
		if tt.email == dummyUser.Email {
			u, err := models.Users.GetByEmail(tt.email)
			if err != nil {
				t.Error("failed to get user: ", err)
			}

			token = u.Token.PlainText
		} else {
			token = tt.token
		}
		authHeader := "Bearer " + token
		_, err := models.Tokens.Authenticate(authHeader)
		if tt.errorExpected && err == nil {
			t.Errorf("%s: %s", tt.name, tt.message)
		} else if !tt.errorExpected && err != nil {
			t.Errorf("%s: %s - %s", tt.name, tt.message, err)
		} else {
			t.Logf("%s", tt.name)
		}
	}
}

func TestToken_Delete(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	err = models.Tokens.Delete(u.Token.ID)
	if err != nil {
		t.Error("error deleting token: ", err)
	}
}

func TestToken_DeleteByToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	err = models.Tokens.DeleteByToken(u.Token.PlainText)
	if err != nil {
		t.Error("error deleting token: ", err)
	}
}

func TestToken_ExpiredToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	token, err := models.Tokens.GenerateToken(u.ID, -time.Hour)
	if err != nil {
		t.Error("err generating token: ", err)
	}

	err = models.Tokens.Insert(*token, *u)
	if err != nil {
		t.Error("failed to insert the token")
	}

	valid, err := models.Tokens.Validate(token.PlainText)
	if err == nil {
		t.Error("failed to validate the token", err)
	}
	if valid {
		t.Error("expired token was validated")
	}
}
