package middlewares

import (
	"net/http"

	"github.com/proGabby/simple_auth_todo_api/pkg/models"
	"github.com/proGabby/simple_auth_todo_api/pkg/utils"
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
func (m *PermissionMiddleware) Authorize(permittedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context
		user, ok := r.Context().Value("user").(*models.User)
		if !ok || user == nil {
			utils.HandleError(map[string]interface{}{"error": "Unauthorized"}, http.StatusUnauthorized, w)
			return
		}

		// Check if the user has the required permissions
		if !m.hasPermission(user, r, permittedRoles) {
			utils.HandleError(map[string]interface{}{"error": "Forbidden", "message": "You are not permitted"}, http.StatusForbidden, w)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

func (m *PermissionMiddleware) hasPermission(user *models.User, r *http.Request, permittedRoles []string) bool {
	// Check if the user's role is in the provided roles slice
	for _, role := range permittedRoles {
		if user.Role == role {
			return true
		}
	}

	return false
}
