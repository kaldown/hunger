package user

import (
	"errors"
	"net/http"
)

type User struct {
	Username string
}

type userKey string

func (u *User) GetUser(r *http.Request) (*User, error) {
	ctx := r.Context()
	if ctx == nil {
		return nil, errors.New("No context")
	}

	u, ok := ctx.Value(userKey("user")).(*User)
	if !ok {
		return nil, errors.New("Can not parse user")
	}
	return u, nil
}
