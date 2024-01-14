package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codewithtoucans/goweb/rand"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Session struct {
	ID        int
	userID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *pgxpool.Pool
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := rand.String(rand.SessionTokenBytes)
	if err != nil {
		return nil, fmt.Errorf("create %w", err)
	}
	session := Session{userID: userID, Token: token, TokenHash: hash(token)}
	queryRow := ss.DB.QueryRow(context.Background(), `update sessions set token_hash=$2 where user_id=$1 returning id`, session.userID, session.TokenHash)
	err = queryRow.Scan(&session.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		row := ss.DB.QueryRow(context.Background(), `insert into sessions (user_id, token_hash) values($1, $2) returning id`, userID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		log.Printf("get session from db was error")
		return nil, fmt.Errorf("query session was error %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := hash(token)
	var user User
	row := ss.DB.QueryRow(context.Background(), `select users.id, users.email, users.password_hash from users join sessions on sessions.user_id = users.id where sessions.token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		log.Printf("query user was error")
		return nil, fmt.Errorf("query user error %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(value string) error {
	tokenHash := hash(value)
	_, err := ss.DB.Exec(context.Background(), `delete from sessions where token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete session was error %w", err)
	}
	return nil
}

func hash(token string) string {
	sum256 := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(sum256[:])
}
