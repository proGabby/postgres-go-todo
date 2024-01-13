package models

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserStore is responsible for interacting with the user data in the database.
type UserStore struct {
	DB *sql.DB
}

// NewUserStore creates a new UserStore instance.
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

// VerifyUserCredentials verifies the user's credentials and returns the user.
func (us *UserStore) VerifyUserCredentials(username, password string) (*User, error) {

	var user User
	query := "SELECT id, username, password, role FROM users WHERE username = $1"
	err := us.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func (us *UserStore) GetUserByID(userID int) (*User, error) {
	var user User
	query := "SELECT id, username, role FROM users WHERE id = $1"
	err := us.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user in the database.
func (us *UserStore) CreateUser(username, password, role string) (*User, error) {

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if er != nil {
		return nil, er
	}

	var userID int
	query := "INSERT INTO users(username, password, role) VALUES($1, $2, $3) RETURNING id"
	err := us.DB.QueryRow(query, username, hashedPassword, role).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	createdUser := &User{
		ID:       userID,
		Username: username,
		Role:     role,
	}

	return createdUser, nil
}

// UpdateUser updates an existing user in the database.
func (us *UserStore) UpdateUser(userID int, username, password, role string) (*User, error) {
	var updatedUserID int
	query := "UPDATE users SET username = $2, password = $3, role = $4 WHERE id = $1 RETURNING id"
	err := us.DB.QueryRow(query, userID, username, password, role).Scan(&updatedUserID)
	if err != nil {
		return nil, err
	}

	updatedUser := &User{
		ID:       updatedUserID,
		Username: username,
		Role:     role,
	}

	return updatedUser, nil
}

// DeleteUser deletes a user from the database.
func (us *UserStore) DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := us.DB.Exec(query, userID)
	return err
}
