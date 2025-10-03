package repository

import (
	"golang-redis/internal/entity"
	"sync"
)

type SessionRepository interface {
	Save(session entity.Session) error
	GetByToken(token string) (*entity.Session, error)
	Delete(token string) error
}

type sessionRepo struct {
	sessions map[string]entity.Session
	mu       sync.Mutex
}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{
		sessions: make(map[string]entity.Session),
	}
}

func (r *sessionRepo) Save(session entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.Token] = session
	return nil
}

func (r *sessionRepo) GetByToken(token string) (*entity.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[token]; ok {
		return &s, nil
	}
	return nil, nil
}

func (r *sessionRepo) Delete(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, token)
	return nil
}
