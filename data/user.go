package data

import (
	"errors"
	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	IsActive  int       `db:"is_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll() ([]*User, error) {
	var all []*User
	collection := upper.Collection(u.Table())
	res := collection.Find().OrderBy("created_at DESC")

	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var user *User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email": email})

	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	token, err := GetTokenByUserId(user.ID)
	if err != nil {
		if !errors.Is(err, up.ErrNilRecord) && !errors.Is(err, up.ErrNoMoreRows) {
			return nil, err
		}
	}

	user.Token = token

	return user, nil
}

func (u *User) Get(id int) (*User, error) {
	var user *User
	collection := upper.Collection(u.Table())
	res := collection.Find(id)

	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	token, err := GetTokenByUserId(user.ID)
	if err != nil {
		if !errors.Is(err, up.ErrNilRecord) && !errors.Is(err, up.ErrNoMoreRows) {
			return nil, err
		}
	}

	user.Token = token

	return user, nil
}

func (u *User) Update(user User) error {
	user.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(user.ID)

	err := res.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Insert(user User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	user.Password = string(hash)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(user)

	if err != nil {
		return 0, err
	}
	id := getInsertId(res.ID())

	return id, nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)

	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ResetPassword(id int, pass string) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)

	user, err := u.Get(id)
	if err != nil {
		return err
	}

	user.Password = string(newHash)
	user.UpdatedAt = time.Now()

	err = user.Update(*u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(pass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
