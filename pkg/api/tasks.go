package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

// tasksHandler handles GET /api/tasks to retrieve all tasks from the scheduler table
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all tasks from the database, ordered by date, replacing NULL with empty strings
	var tasks []Task
	err := db.DB.Select(&tasks, `SELECT id, date, title, COALESCE(comment, '') AS comment, COALESCE(repeat, '') AS repeat FROM scheduler ORDER BY date`)
	if err != nil {
		writeError(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	// Ensure tasks is not nil
	if tasks == nil {
		tasks = []Task{}
	}

	// Return tasks as JSON in the format {"tasks": [...]}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Task{"tasks": tasks})
}
