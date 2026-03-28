package lib

import (
	"encoding/json"
	"sort"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
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

//export listTodos
func listTodos(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	keys, err := db.List("todo/")
	if err != nil {
		return respondError(h, 500, "failed to list todos")
	}
	sort.Strings(keys)

	todos := make([]Todo, 0, len(keys))
	for _, key := range keys {
		raw, getErr := db.Get(key)
		if getErr != nil {
			continue
		}

		var todo Todo
		if unmarshalErr := json.Unmarshal(raw, &todo); unmarshalErr != nil {
			continue
		}
		todos = append(todos, todo)
	}

	return respondJSON(h, 200, todos)
}

func respondJSON(h event.HTTP, status int, payload interface{}) uint32 {
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

func respondError(h event.HTTP, status int, message string) uint32 {
	type errorResponse struct {
		Error string `json:"error"`
	}
	return respondJSON(h, status, errorResponse{Error: message})
}
