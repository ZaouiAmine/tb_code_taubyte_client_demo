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

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

const dbMatch = "appdata"

//export createNote
func createNote(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	body, err := io.ReadAll(h.Body())
	if err != nil {
		return respondError(h, 400, "failed to read request body")
	}
	defer h.Body().Close()

	var note Note
	if err := json.Unmarshal(body, &note); err != nil {
		return respondError(h, 400, "invalid note payload")
	}
	if note.Title == "" {
		return respondError(h, 400, "title is required")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	note.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	note.CreatedAt = now
	note.UpdatedAt = now

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	encoded, err := json.Marshal(note)
	if err != nil {
		return respondError(h, 500, "failed to encode note")
	}
	if err := db.Put("note/"+note.ID, encoded); err != nil {
		return respondError(h, 500, "failed to store note")
	}

	return respondJSON(h, 201, note)
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
