package main

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}

func (user *User) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return PrivateKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token.Valid false")
	}

	return nil
}

type UsersCollection struct {
	mu    sync.RWMutex
	users []*User
}

func (uc *UsersCollection) AddUser(nickname string) (*User, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    "server tic tac toe",
		Subject:   nickname,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(PrivateKey)
	if err != nil {
		return nil, err
	}

	user := User{Nickname: nickname, Token: signedToken}

	uc.mu.Lock()
	uc.users = append(uc.users, &user)
	uc.mu.Unlock()

	return &user, nil
}

func (uc *UsersCollection) GetUserByToken(token string) (*User, error) {
	uc.mu.RLock()
	idx := slices.IndexFunc(uc.users, func(u *User) bool {
		return u.Token == token
	})
	uc.mu.RUnlock()

	if idx == -1 {
		return nil, errors.New("User not found")
	}

	uc.mu.RLock()
	user := uc.users[idx]
	uc.mu.RUnlock()

	return user, nil
}
