package middlewares

import (
	"net/http"
	"strings"
	"github.com/proGabby/simple_auth_todo_api/pkg/models"
)

// PermissionMiddleware handles user authorization based on roles.
type PermissionMiddleware struct {
	AuthMiddleware *AuthMiddleware
	TodoStore      models.TodoStore
}

func NewPermissionMiddleware(authMiddleware *AuthMiddleware, todoStore models.TodoStore) *PermissionMiddleware {
	return &PermissionMiddleware{AuthMiddleware: authMiddleware, TodoStore: todoStore}
}

// Authorize is the middleware function that checks if the user has the required permissions.
func (m *PermissionMiddleware) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context
		user, ok := r.Context().Value("user").(*models.User)
		if !ok || user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user has the required permissions
		if !m.hasPermission(user, r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

func (m *PermissionMiddleware) hasPermission(user *models.User, r *http.Request) bool {
	// Extract required permission from the route (for simplicity, assuming it's in the path)
	requiredPermission := extractRequiredPermission(r)

	// Check if the user has the required permission based on their role
	switch user.Role {
	case "admin":
		// Admin has all permissions
		return true
	case "user":
		// User has limited permissions (for demonstration purposes)
		return requiredPermission != "admin"
	default:
		return false
	}
}

func extractRequiredPermission(r *http.Request) string {
	// For simplicity, assuming the required permission is the last part of the route path
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) > 0 {
		return pathSegments[len(pathSegments)-1]
	}
	return ""
}