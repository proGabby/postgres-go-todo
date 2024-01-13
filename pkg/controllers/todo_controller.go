package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/proGabby/simple_auth_todo_api/pkg/models"
)

// TodoController handles todo-related HTTP requests.
type TodoController struct {
	TodoStore models.TodoStore
}

// NewTodoController creates a new TodoController instance.
func NewTodoController(todoStore models.TodoStore) *TodoController {
	return &TodoController{TodoStore: todoStore}
}

// GetTodosByUser retrieves all todos for a specific user.
func (c *TodoController) GetTodosByUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (assuming it was set during authentication)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve todos for the user
	todos, err := c.TodoStore.GetTodosByUserID(user.ID)
	if err != nil {
		http.Error(w, "Error retrieving todos", http.StatusInternalServerError)
		return
	}

	// Return todos in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo creates a new todo for the authenticated user.
func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (assuming it was set during authentication)
	user, ok := r.Context().Value("user").(*models.User)

	if !ok || user == nil {
		jsonResponse := map[string]interface{}{
			"error": "Unauthorized",
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	// Parse the JSON request body
	var newTodo models.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create the todo
	createdTodo, err := c.TodoStore.CreateTodo(user.ID, newTodo.Title, "active")
	if err != nil {
		http.Error(w, "Error creating todo", http.StatusInternalServerError)
		return
	}

	// Return the created todo in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTodo)
}

// UpdateTodo updates an existing todo for the authenticated user.
// func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
// 	// Extract user ID from the request context (assuming it was set during authentication)
// 	user, ok := r.Context().Value("user").(*models.User)
// 	if !ok || user == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse todo ID from the request URL
// 	todoID, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil {
// 		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse the JSON request body
// 	var updatedTodo models.Todo
// 	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Update the todo
// 	updatedTodo, err = c.TodoStore.UpdateTodo(todoID, updatedTodo.Title, updatedTodo.Status)
// 	if err != nil {
// 		http.Error(w, "Error updating todo", http.StatusInternalServerError)
// 		return
// 	}

// 	// Return the updated todo in the response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(updatedTodo)
// }

// DeleteTodo deletes an existing todo for the authenticated user.
func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (assuming it was set during authentication)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse todo ID from the request URL
	todoID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Delete the todo
	err = c.TodoStore.DeleteTodo(todoID)
	if err != nil {
		http.Error(w, "Error deleting todo", http.StatusInternalServerError)
		return
	}

	// Return success in the response
	w.WriteHeader(http.StatusOK)
}
