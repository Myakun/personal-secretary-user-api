package jwtrefreshtoken

import (
	"time"
)

type JwtRefreshToken struct {
	createdAt  time.Time
	id         string
	isInserted bool
	userId     string
	token      string
}

func NewJwtRefreshToken(createdAt time.Time, token string, userId string) *JwtRefreshToken {
	return &JwtRefreshToken{
		createdAt:  createdAt,
		isInserted: false,
		token:      token,
		userId:     userId,
	}
}

func (entity *JwtRefreshToken) GetCreatedAt() time.Time {
	return entity.createdAt
}

func (entity *JwtRefreshToken) GetId() string {
	return entity.id
}

func (entity *JwtRefreshToken) GetUserId() string {
	return entity.userId
}

func (entity *JwtRefreshToken) GetToken() string {
	return entity.token
}

func (entity *JwtRefreshToken) IsInserted() bool {
	return entity.isInserted
}

func (entity *JwtRefreshToken) setIsInserted(isInserted bool) {
	entity.isInserted = isInserted
}
