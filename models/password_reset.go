package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/codewithtoucans/goweb/rand"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DefaultResetDuration = 1 * time.Hour

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB       *pgxpool.Pool
	Duration time.Duration
}

func (p *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userID int
	row := p.DB.QueryRow(context.Background(), `select id from users where email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create %w", err)
	}
	token, err := rand.String(rand.SessionTokenBytes)
	if err != nil {
		return nil, fmt.Errorf("create %w", err)
	}
	duration := p.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := PasswordReset{UserID: userID, Token: token, TokenHash: p.hash(token), ExpiresAt: time.Now().Add(duration)}
	row = p.DB.QueryRow(context.Background(), `insert into password_resets (user_id, token_hash, expires_at) values ($1, $2, $3) on conflict (user_id) do update set token_hash = $2, expires_at = $3 returning id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &pwReset, nil
}

func (p *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := p.hash(token)
	var user User
	var pwReset PasswordReset
	row := p.DB.QueryRow(context.Background(), `
		select password_resets.id,
		password_resets.expires_at,
		users.id,
		users.email,
		users.password_hash
		from password_resets
		join users on users.id = password_resets.user_id
		where password_resets.token_hash = $1
	`, tokenHash)
	err := row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = p.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume password_resets error %w", err)
	}
	return &user, nil
}

func (p *PasswordResetService) hash(token string) string {
	sum256 := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(sum256[:])
}

func (p *PasswordResetService) delete(id int) error {
	_, err := p.DB.Exec(context.Background(), `delete from password_resets where id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete password_resets error %w", err)
	}
	return nil
}
