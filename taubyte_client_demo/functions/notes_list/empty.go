package lib

import (
	"encoding/json"
	"sort"
	"strings"

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

//export listNotes
func listNotes(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	search, _ := h.Query().Get("q")
	search = strings.ToLower(strings.TrimSpace(search))

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	keys, err := db.List("note/")
	if err != nil {
		return respondError(h, 500, "failed to list notes")
	}
	sort.Strings(keys)

	notes := make([]Note, 0, len(keys))
	for _, key := range keys {
		raw, getErr := db.Get(key)
		if getErr != nil {
			continue
		}

		var note Note
		if unmarshalErr := json.Unmarshal(raw, &note); unmarshalErr != nil {
			continue
		}

		if search != "" {
			fullText := strings.ToLower(note.Title + " " + note.Content + " " + strings.Join(note.Tags, " "))
			if !strings.Contains(fullText, search) {
				continue
			}
		}

		notes = append(notes, note)
	}

	return respondJSON(h, 200, notes)
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
