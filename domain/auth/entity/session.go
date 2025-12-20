package entity

import (
	"time"

	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
)

type Session struct {
	trait.Createable
	trait.Updateable

	ID           string
	UserId       string
	AccessToken  string
	RefreshToken string
	RevokedAt    *time.Time
}

func NewSession(userId, accessToken, refreshToken string) *Session {
	session := &Session{
		ID:           util.GenerateUUID(),
		UserId:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	session.MarkCreate()

	return session
}

func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *Session) Revoke() {
	now := time.Now().UTC()
	s.RevokedAt = &now

	s.MarkUpdate()
}

func (s *Session) IsAccessTokenEqual(accessToken string) bool {
	return s.AccessToken == accessToken
}

func (s *Session) IsRefreshTokenEqual(refreshToken string) bool {
	return s.RefreshToken == refreshToken
}
