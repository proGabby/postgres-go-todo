package utils

import (
	"encoding/json"
	"net/http"
)

func HandleError( jsonResponse map[string]interface{}, httpStatusCode int, w http.ResponseWriter) {
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatusCode)

		if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
			// Handle encoding error if necessary
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

}
