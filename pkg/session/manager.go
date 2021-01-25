package session

import (
	"net/http"
	"time"
)

type Manager interface {
	New() (string, error)
	Check(r *http.Request) (bool, error)
	Remove(expireAfter time.Duration) ([]string, error)
}
