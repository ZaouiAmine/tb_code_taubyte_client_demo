package lib

import (
	"encoding/json"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	httpevent "github.com/taubyte/go-sdk/http/event"
)

const dbMatch = "appdata"

//export deleteNote
func deleteNote(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	id, _ := h.Query().Get("id")
	if id == "" {
		return respondError(h, 400, "id query parameter is required")
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	key := "note/" + id
	if _, err := db.Get(key); err != nil {
		return respondError(h, 404, "note not found")
	}
	if err := db.Delete(key); err != nil {
		return respondError(h, 500, "failed to delete note")
	}

	type response struct {
		Deleted bool   `json:"deleted"`
		ID      string `json:"id"`
	}
	return respondJSON(h, 200, response{Deleted: true, ID: id})
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
