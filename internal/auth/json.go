package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	a "github.com/juztin/statictls/pkg/auth"
)

type JsonAuthenticator struct {
	usersPath    string
	users        map[string]string
	lastModified time.Time
}

func (j *JsonAuthenticator) loadUsers() error {
	stats, err := os.Stat(j.usersPath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(j.usersPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &j.users)
	if err == nil {
		j.lastModified = stats.ModTime()
	}
	return err
}

func (j *JsonAuthenticator) Authenticate(username, password string) error {
	stats, err := os.Stat(j.usersPath)
	if err != nil {
		return err
	}
	if stats.ModTime().Sub(j.lastModified) > 0 {
		j.loadUsers()
	}
	hashed, ok := j.users[username]
	if !ok {
		return a.ErrInvalidCredentials
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func NewJson(usersPath string) *JsonAuthenticator {
	return &JsonAuthenticator{usersPath: usersPath}
}
