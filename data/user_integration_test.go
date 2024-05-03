//go:build integration

// run tests with this command: go test --cover . -v --tags integration --count=1
package data_test

import (
	"myapp/data"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testUser = data.User{
	FirstName: "Some",
	LastName:  "Guy",
	Email:     "me@here.com",
	IsActive:  1,
	Password:  "password",
}

func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("wrong table name returned: ", s)
	}
}

func TestUser_Insert(t *testing.T) {
	id, err := models.Users.Insert(testUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	if id == 0 {
		t.Error("0 returned as id after insert")
	}
}

func TestUser_Get(t *testing.T) {
	u, err := models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("id of returned user is 0: ", err)
	}
}

func TestUser_GetAll(t *testing.T) {
	_, err := models.Users.GetAll()
	if err != nil {
		t.Error("failed to get user: ", err)
	}
}

func TestUser_GetByEmail(t *testing.T) {
	u, err := models.Users.GetByEmail(testUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("id of returned user is 0: ", err)
	}
}

func TestUser_Update(t *testing.T) {
	u, err := models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	u.LastName = "Smith"
	err = u.Update(*u)
	if err != nil {
		t.Error("failed to update user: ", err)
	}

	u, err = models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.LastName != "Smith" {
		t.Error("last name not updated in database")
	}
}

func TestUser_PasswordMatches(t *testing.T) {
	u, err := models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	matches, err := u.PasswordMatches(testUser.Password)
	if err != nil {
		t.Error("error checkign match: ", err)
	}
	if !matches {
		t.Error("user password does not match")
	}

	matches, err = u.PasswordMatches("worng password")
	if err != nil {
		t.Error("error checkign match: ", err)
	}
	if matches {
		t.Error("user password match when it should not!")
	}
}

func TestUser_ResetPassword(t *testing.T) {
	newPassword := "new_password"
	err := models.Users.ResetPassword(1, newPassword)
	if err != nil {
		t.Error("error reseting password: ", err)
	}

	err = models.Users.ResetPassword(20000, newPassword)
	if err == nil {
		t.Error("did not get an error reseting password for none-existing user", err)
	}
}

func TestUser_Delete(t *testing.T) {
	err := models.Users.Delete(1)
	if err != nil {
		t.Error("failed to delete user: ", err)
	}

	_, err = models.Users.Get(1)
	if err == nil {
		t.Error("retrived user who supposed to be deleted!")
	}
}
