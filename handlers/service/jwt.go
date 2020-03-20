package service

import (
	"errors"

	"github.com/google/uuid"
)

var storage = make(map[string]*Token)

type StoredToken interface {
	Store(string) *Token
	FindByKey(string) (string, error)
}

type Token struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

func NewStoredToken(token string) *Token {
	return &Token{
		Token: token,
	}
}

func (t *Token) Store() *Token {
	t.Key = uuid.New().String()

	storage[t.Key] = t

	return t
}

func (t Token) FindByKey(key string) (*Token, error) {
	if t, ok := storage[key]; ok {
		return t, nil
	}

	return nil, errors.New("Missing token in storage")
}
