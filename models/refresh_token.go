package models

import (
	"time"

	"github.com/BatJoz21/my-online-shop-go-api/database"
)

type RefreshToken struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	DeviceName *string    `json:"device_name"`
	TokenHash  string     `json:"token_hash"`
	ExpiresAt  *time.Time `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  *time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func GetRefreshTokenByHashedToken(hashed string) (*RefreshToken, error) {
	query := `SELECT
		id,
		user_id,
		device_name,
		token_hash,
		expires_at,
		revoked_at,
		created_at
	FROM refresh_tokens WHERE token_hash = ?`
	row := database.DB.QueryRow(query, hashed)

	var refreshToken RefreshToken
	err := row.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.DeviceName, &refreshToken.TokenHash, &refreshToken.ExpiresAt, &refreshToken.RevokedAt, &refreshToken.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshToken) Save() error {
	query := `INSERT INTO refresh_tokens(user_id, device_name, token_hash, expires_at)
		VALUES (?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.UserID, r.DeviceName, r.TokenHash, r.ExpiresAt)
	if err != nil {
		return err
	}

	tokenId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = tokenId

	return nil
}

func (r *RefreshToken) Revoke() error {
	query := `UPDATE refresh_tokens SET revoked_at = ? WHERE token_hash = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(time.Now(), r.TokenHash)
	if err != nil {
		return err
	}

	return nil
}
