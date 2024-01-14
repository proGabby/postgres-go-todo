package models

import (
	"database/sql"
	"fmt"
)

// Todo represents a task in the system.
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"user_id"`
}

// TodoStore is responsible for interacting with the todo data in the database.
type TodoStore struct {
	DB *sql.DB
}

// NewTodoStore creates a new TodoStore instance.
func NewTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{DB: db}
}

// GetTodosByUserID retrieves all todos for a given user ID.
func (ts *TodoStore) GetTodosByUserID(userID int) ([]Todo, error) {
	var todos []Todo
	query := "SELECT id, title, status FROM todos WHERE user_id = $1"
	rows, err := ts.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			return nil, err
		}
		todo.UserID = userID
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodoByID retrieves a todo by its ID.
func (ts *TodoStore) GetTodoByID(todoID int) (*Todo, error) {
	var todo Todo
	query := "SELECT id, title, status, user_id FROM todos WHERE id = $1"
	err := ts.DB.QueryRow(query, todoID).Scan(&todo.ID, &todo.Title, &todo.Status, &todo.UserID)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// CreateTodo creates a new todo in the database.
func (ts *TodoStore) CreateTodo(userID int, title, status string) (*Todo, error) {
	var todoID int
	query := "INSERT INTO todos(title, status, user_id) VALUES($1, $2, $3) RETURNING id"
	err := ts.DB.QueryRow(query, title, status, userID).Scan(&todoID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	createdTodo := &Todo{
		ID:     todoID,
		Title:  title,
		Status: status,
		UserID: userID,
	}

	return createdTodo, nil
}

// UpdateTodo updates an existing todo in the database.
func (ts *TodoStore) UpdateTodo(todoID int, title, status string, userId int) (*Todo, error) {
	var updatedTodo Todo
	query := "UPDATE todos SET title = COALESCE(NULLIF($2, ''), title), status = COALESCE(NULLIF($3, ''), status) WHERE id = $1 RETURNING id, title, status"
	err := ts.DB.QueryRow(query, todoID, title, status).Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Status)
	if err != nil {
		return nil, err
	}

	newUpdatedTodo := &Todo{
		ID:     updatedTodo.ID,
		Title:  updatedTodo.Title,
		Status: updatedTodo.Status,
		UserID: userId,
	}

	return newUpdatedTodo, nil
}

// DeleteTodo deletes a todo from the database.
func (ts *TodoStore) DeleteTodo(todoID int) error {
	query := "DELETE FROM todos WHERE id = $1"
	_, err := ts.DB.Exec(query, todoID)
	return err
}
