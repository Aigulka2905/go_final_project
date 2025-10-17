package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Create a new task
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			writeError(w, "title is required", http.StatusBadRequest)
			return
		}

		// Set date to today if empty
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if task.Date == "" {
			task.Date = today.Format(dateFormat)
		}

		// Validate date format
		date, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			writeError(w, "invalid date format", http.StatusBadRequest)
			return
		}

		// If date is before today, set it to today
		if date.Before(today) {
			task.Date = today.Format(dateFormat)
		}

		// Calculate next date if repeat is set and date is not today with repeat="d 1"
		if task.Repeat != "" && !(task.Repeat == "d 1" && task.Date == today.Format(dateFormat)) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = next
		}

		// Insert task into database
		result, err := db.DB.NamedExec(
			`INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`,
			task,
		)
		if err != nil {
			writeError(w, "failed to save task", http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			writeError(w, "failed to get task ID", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("%d", id)})

	case http.MethodPut:
		// Update an existing task
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeError(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			writeError(w, "title is required", http.StatusBadRequest)
			return
		}

		// Set date to today if empty
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if task.Date == "" {
			task.Date = today.Format(dateFormat)
		}

		// Validate date format
		date, err := time.Parse(dateFormat, task.Date)
		if err != nil {
			writeError(w, "invalid date format", http.StatusBadRequest)
			return
		}

		// If date is before today, set it to today
		if date.Before(today) {
			task.Date = today.Format(dateFormat)
		}

		// Calculate next date if repeat is set and date is not today with repeat="d 1"
		if task.Repeat != "" && !(task.Repeat == "d 1" && task.Date == today.Format(dateFormat)) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = next
		}

		// Update task in database
		result, err := db.DB.NamedExec(
			`UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id`,
			task,
		)
		if err != nil {
			writeError(w, "failed to update task", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeError(w, "failed to check update result", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			writeError(w, "task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})

	case http.MethodDelete:
		// Delete a task
		var req DeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			writeError(w, "id is required", http.StatusBadRequest)
			return
		}

		// Delete task from database
		result, err := db.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, req.ID)
		if err != nil {
			writeError(w, "failed to delete task", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeError(w, "failed to check delete result", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			writeError(w, "task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})

	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
