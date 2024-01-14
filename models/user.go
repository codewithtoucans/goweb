package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/codewithtoucans/goweb/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *pgxpool.Pool
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user generate password error: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{Email: email, PasswordHash: passwordHash}
	row := us.DB.QueryRow(context.Background(), `INSERT INTO users (email, password_hash) values ($1, $2) RETURNING id`, email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{Email: email}
	row := us.DB.QueryRow(context.Background(), `select id, password_hash from users where email = $1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate user error %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("authenticate hashing: %v\n", password)
		return nil, errors.NewPublicError(ErrNotFound, ErrNotFound.Error())
	}
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("update password error %s\n", err.Error())
		return fmt.Errorf("update password error %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(context.Background(), `update users set password_hash=$2 where id=$1`, userID, passwordHash)
	if err != nil {
		log.Printf("update password error %s\n", err.Error())
		return fmt.Errorf("update password error %w", err)
	}
	return nil
}

// check user exists
func (us *UserService) CheckUserExist(email string) (bool, error) {
	var exists bool
	err := us.DB.QueryRow(context.Background(), `select exists (select 1 from users where email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exist error %w", err)
	}
	return exists, nil
}
