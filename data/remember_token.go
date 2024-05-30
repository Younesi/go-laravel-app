package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type RememberToken struct {
	ID            int       `db:"id,omitempty"`
	UserId        int       `db:"user_id"`
	RememberToken string    `db:"remember_token"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *RememberToken) Table() string {
	return "remember_tokens"
}

func (r *RememberToken) InsertToken(userId int, token string) error {
	collection := upper.Collection(r.Table())
	rememberToken := RememberToken{
		UserId:        userId,
		RememberToken: token,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := collection.Insert(rememberToken)
	if err != nil {
		return err
	}

	return nil
}

func (r *RememberToken) Delete(rt string) error {
	collection := upper.Collection(r.Table())
	res := collection.Find(up.Cond{"remember_token": rt})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}
