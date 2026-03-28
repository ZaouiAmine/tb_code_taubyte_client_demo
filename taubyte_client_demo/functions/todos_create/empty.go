package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	httpevent "github.com/taubyte/go-sdk/http/event"
)

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

const dbMatch = "appdata"

//export createTodo
func createTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	body, err := io.ReadAll(h.Body())
	if err != nil {
		return respondError(h, 400, "failed to read request body")
	}
	defer h.Body().Close()

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		return respondError(h, 400, "invalid todo payload")
	}
	if todo.Title == "" {
		return respondError(h, 400, "title is required")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	todo.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	todo.CreatedAt = now
	todo.UpdatedAt = now

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	encoded, err := json.Marshal(todo)
	if err != nil {
		return respondError(h, 500, "failed to encode todo")
	}
	if err := db.Put("todo/"+todo.ID, encoded); err != nil {
		return respondError(h, 500, "failed to store todo")
	}

	return respondJSON(h, 201, todo)
}

func respondJSON(h httpevent.Event, status int, payload interface{}) uint32 {
	body, err := json.Marshal(payload)
	if err != nil {
		h.Return(500)
		return 0
	}
	h.Headers().Set("Content-Type", "application/json")
	h.Write(body)
	h.Return(status)
	return 0
}

func respondError(h httpevent.Event, status int, message string) uint32 {
	type errorResponse struct {
		Error string `json:"error"`
	}
	return respondJSON(h, status, errorResponse{Error: message})
}
