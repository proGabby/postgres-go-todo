package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/proGabby/simple_auth_todo_api/pkg/models"
	"github.com/proGabby/simple_auth_todo_api/pkg/utils"
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
	// Extract user ID from the request context that was set during authentication
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		utils.HandleError(map[string]interface{}{
			"error": "Unauthorized",
		}, http.StatusUnauthorized, w)
		return
	}

	// Retrieve todos for the user
	todos, err := c.TodoStore.GetTodosByUserID(user.ID)
	if err != nil {
		fmt.Print(err)
		utils.HandleError(map[string]interface{}{
			"error":   "internal server error",
			"message": "error retrieving todos",
		}, http.StatusInternalServerError, w)
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

		utils.HandleError(jsonResponse, http.StatusUnauthorized, w)

		return
	}

	// Parse the JSON request body
	var newTodo models.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid request body",
		}, http.StatusBadRequest, w)
		return
	}

	// Create the todo
	createdTodo, err := c.TodoStore.CreateTodo(user.ID, newTodo.Title, "active")
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Internal Server Error",
			"message": "Error creating todo",
		}, http.StatusInternalServerError, w)
		return
	}

	// Return the created todo in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTodo)
}

// UpdateTodo updates an existing todo for the authenticated user.
func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context (assuming it was set during authentication)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		utils.HandleError(map[string]interface{}{
			"error": "Unauthorized",
		}, http.StatusUnauthorized, w)
		return
	}

	// Parse todo ID from the request URL
	todoID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid todo ID",
		}, http.StatusBadRequest, w)
		return
	}

	// Parse the JSON request body
	var updatedTodo models.Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid request body",
		}, http.StatusBadRequest, w)
		return
	}

	// Update the todo
	newUpdatedTodo, err := c.TodoStore.UpdateTodo(todoID, updatedTodo.Title, updatedTodo.Status, user.ID)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Internal Server Error",
			"message": "Error updating todo",
		}, http.StatusInternalServerError, w)
		return
	}

	// Return the updated todo in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUpdatedTodo)
}

func (c *TodoController) GetSingleTodo(w http.ResponseWriter, r *http.Request) {

	// Extract user ID from the request context (assuming it was set during authentication)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		jsonResponse := map[string]interface{}{
			"error": "Unauthorized",
		}

		utils.HandleError(jsonResponse, http.StatusUnauthorized, w)

		return
	}

	// Parse todo ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	todoID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(todoID)
		fmt.Println(err)
		utils.HandleError(map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid todo ID",
		}, http.StatusBadRequest, w)
		return
	}

	// Retrieve the todo
	todo, err := c.TodoStore.GetTodoByID(todoID)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Data Error",
			"message": "Error retrieving todo",
		}, http.StatusInternalServerError, w)
		return
	}

	// Return the todo in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes an existing todo for the authenticated user.
func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request context that was set during authentication
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		utils.HandleError(map[string]interface{}{
			"error": "Unauthorized",
		}, http.StatusUnauthorized, w)
		return
	}

	// Parse todo ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]
	//r.URL.Query().Get("id")

	todoID, err := strconv.Atoi(id)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Bad Request",
			"message": "Invalid todo ID",
		}, http.StatusBadRequest, w)
		return
	}

	// Delete the todo
	err = c.TodoStore.DeleteTodo(todoID)
	if err != nil {
		utils.HandleError(map[string]interface{}{
			"error":   "Data Error",
			"message": "Error deleting todo",
		}, http.StatusInternalServerError, w)
		return
	}

	// Return success in the response
	w.WriteHeader(http.StatusOK)
}
