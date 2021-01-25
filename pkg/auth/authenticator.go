package auth

import "errors"

var ErrInvalidCredentials = errors.New("auth: invalid username")

type Authenticator interface {
	Authenticate(username, password string) error
}
