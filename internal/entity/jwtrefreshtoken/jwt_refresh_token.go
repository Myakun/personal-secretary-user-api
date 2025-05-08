package jwtrefreshtoken

import (
	"time"
)

type RefreshToken struct {
	createdAt time.Time
	id        string
	userId    string
	token     string
}

func NewRefreshToken(userId, token string) *RefreshToken {
	return &RefreshToken{
		userId: userId,
		token:  token,
	}
}

func (r *RefreshToken) GetCreatedAt() time.Time {
	return r.createdAt
}

func (r *RefreshToken) GetId() string {
	return r.id
}

func (r *RefreshToken) GetUserId() string {
	return r.userId
}

func (r *RefreshToken) GetToken() string {
	return r.token
}
