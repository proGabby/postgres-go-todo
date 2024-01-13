package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/proGabby/simple_auth_todo_api/pkg/middlewares"
	"github.com/proGabby/simple_auth_todo_api/pkg/models"
)

// UserController handles user-related HTTP requests.
type UserController struct {
	UserStore      models.UserStore
	authMiddleware middlewares.AuthMiddleware
}

// NewUserController creates a new UserController instance.
func NewUserController(userStore models.UserStore) *UserController {
	return &UserController{UserStore: userStore}
}

// RegisterUser handles user registration.
func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username and password (add more validation as needed)
	if newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Create the user (including password hashing)
	createdUser, err := c.UserStore.CreateUser(newUser.Username, newUser.Password, "user")
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Omit the Password field from the response
	createdUser.Password = ""

	// Return the created user in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}

// LoginUser handles user login.
func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var loginUser models.User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username and password (add more validation as needed)
	if loginUser.Username == "" || loginUser.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Verify user credentials
	user, err := c.UserStore.VerifyUserCredentials(loginUser.Username, loginUser.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token with the user payload
	token, err := c.authMiddleware.GenerateJWTToken(user)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	// Omit the Password field from the response
	user.Password = ""

	// Return the authenticated user and the JWT token in the response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

func (c *UserController) GetUserByToken(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value("user").(*models.User)

	if !ok {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}
	fmt.Print(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
