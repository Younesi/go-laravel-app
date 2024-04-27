package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	up "github.com/upper/db/v4"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	ID        int       `db:"id" json:"id"`
	UserId    int       `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	Email     string    `db:"email" json:"email"`
	PlainText string    `db:"token" json:"token"`
	Hash      []byte    `db:"token_hash" json:"_"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) GetUserForToken(token string) (*User, error) {
	var user *User
	var theToken Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": token})

	err := res.One(&theToken)
	if err != nil {
		return nil, err
	}

	collection = upper.Collection(user.Table())
	res = collection.Find(up.Cond{"id": theToken.UserId})
	err = res.One(&user)
	if err != nil {
		return nil, err
	}

	user.Token = theToken

	return user, nil
}

func (t *Token) GetTokensForUser(id int) ([]*Token, error) {
	var tokens []*Token
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"user_id": id})
	err := res.All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (t *Token) Get(id int) (*Token, error) {
	var token *Token
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Token) GetByToken(plainText string) (*Token, error) {
	var token *Token
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"token": plainText})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *Token) Delete(id int) error {
	collection := upper.Collection(t.Table())

	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteByToken(plainText string) error {
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"token": plainText})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Insert(token Token, u User) error {
	collection := upper.Collection(t.Table())

	// delete existing token
	res := collection.Find(up.Cond{"user_id": u.ID})
	if err := res.Delete(); err != nil {
		return err
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.UserId = u.ID
	token.FirstName = u.FirstName
	token.Email = u.Email
	// token ?

	_, err := collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GenerateToken(userId int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserId:    userId,
		ExpiresAt: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (t *Token) Authenticate(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token wrong size")
	}

	tkn, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no matching token found")
	}

	if tkn.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	user, err := t.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

func (t *Token) Validate(token string) (bool, error) {
	user, err := t.GetUserForToken(token)
	if err != nil {
		return false, errors.New("no matching user found")
	}

	if user.Token.PlainText != "" {
		return false, errors.New("no matching token found")
	}

	if user.Token.ExpiresAt.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}

func GetTokenByUserId(userId int) (Token, error) {
	var token Token
	collection := upper.Collection(token.Table())

	res := collection.Find(up.Cond{"user_id": userId, "expires_at >": time.Now()}).OrderBy("created_at DESC")
	err := res.One(&token)
	if err != nil {
		return token, err
	}

	return token, nil
}
