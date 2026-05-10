package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Authenticate(username, password string) (bool, error) {
	var passwordHash string

	log.Printf("username %s, password %s", username, password)

	err := r.db.QueryRow(
		"SELECT password_hash FROM app_user WHERE username = $1",
		username,
	).Scan(&passwordHash)

	log.Printf("passwordHash %s", passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		fmt.Printf("Authentication failed for user: %s with error %v\n", username, err)
		return false, nil
	}

	fmt.Printf("Authentication successful for user: %s\n", username)
	return true, nil
}

// HashPassword hashes a plain text password for secure storage
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
