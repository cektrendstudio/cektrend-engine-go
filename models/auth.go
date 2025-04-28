package models

import "time"

type CreateTokenRequest struct {
	UserID                int64     `json:"user_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ValidateAccessTokenRequest struct {
	UserID      int64  `json:"user_Id"`
	AccessToken string `json:"access_token"`
}

type ValidateRefreshTokenRequest struct {
	UserID       int64  `json:"user_Id"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateTokenRequest struct {
	TokenID               int64     `json:"token_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type RevokedInactiveTokenRequest struct {
	UserID  int64 `json:"user_id"`
	TokenID int64 `json:"token_id"`
}

type ValidateAccessTokenResponse struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"access_token_expires_at" db:"access_token_expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
}

type ValidateRefreshTokenResponse struct {
	TokenID               int64     `json:"token_id" db:"token_id"`
	UserID                int64     `json:"user_id" db:"user_id"`
	RefreshToken          string    `json:"refresh_token"  db:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at" db:"refresh_token_expires_at"`
}

type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"access_token_expires_at"`
}
