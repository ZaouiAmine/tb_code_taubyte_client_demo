package lib

import (
	"encoding/json"
	"io"
	"time"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
)

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

const dbMatch = "appdata"

//export updateNote
func updateNote(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	id, _ := h.Query().Get("id")
	if id == "" {
		return respondError(h, 400, "id query parameter is required")
	}

	body, err := io.ReadAll(h.Body())
	if err != nil {
		return respondError(h, 400, "failed to read request body")
	}
	defer h.Body().Close()

	var incoming Note
	if err := json.Unmarshal(body, &incoming); err != nil {
		return respondError(h, 400, "invalid note payload")
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	key := "note/" + id
	currentRaw, err := db.Get(key)
	if err != nil {
		return respondError(h, 404, "note not found")
	}

	var current Note
	if err := json.Unmarshal(currentRaw, &current); err != nil {
		return respondError(h, 500, "failed to decode stored note")
	}

	if incoming.Title != "" {
		current.Title = incoming.Title
	}
	if incoming.Content != "" {
		current.Content = incoming.Content
	}
	if len(incoming.Tags) > 0 {
		current.Tags = incoming.Tags
	}
	current.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	encoded, err := json.Marshal(current)
	if err != nil {
		return respondError(h, 500, "failed to encode note")
	}
	if err := db.Put(key, encoded); err != nil {
		return respondError(h, 500, "failed to update note")
	}

	return respondJSON(h, 200, current)
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
