package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

type Memory struct {
	sessions map[string]time.Time
}

func (m *Memory) New() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	session := fmt.Sprintf("%x", b)
	m.sessions[session] = time.Now()
	return session, nil
}

func (m *Memory) Check(r *http.Request) (bool, error) {
	c, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil
		}
		return false, err
	}
	_, ok := m.sessions[c.Value]
	if ok {
		m.sessions[c.Value] = time.Now()
	}
	return ok, nil
}

func (m *Memory) Remove(expireAfter time.Duration) ([]string, error) {
	var cleaned []string
	for k, v := range m.sessions {
		if time.Now().Sub(v) > expireAfter {
			cleaned = append(cleaned, k)
			delete(m.sessions, k)
		}
	}
	return cleaned, nil
}

func NewMemory() *Memory {
	m := &Memory{make(map[string]time.Time)}
	go func() {
		for {
			<-time.After(5 * time.Minute)
			m.Remove(15 * time.Minute)
		}
	}()
	return m
}
