package lib

import (
	"encoding/json"
	"sort"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
)

type ChatMessage struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

const dbMatch = "appdata"

//export listChatMessages
func listChatMessages(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	roomID, _ := h.Query().Get("roomId")
	if roomID == "" {
		return respondError(h, 400, "roomId query parameter is required")
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	prefix := "chatmsg/" + roomID + "/"
	keys, err := db.List(prefix)
	if err != nil {
		return respondError(h, 500, "failed to list messages")
	}
	sort.Strings(keys)

	messages := make([]ChatMessage, 0, len(keys))
	for _, key := range keys {
		raw, getErr := db.Get(key)
		if getErr != nil {
			continue
		}

		var message ChatMessage
		if unmarshalErr := json.Unmarshal(raw, &message); unmarshalErr != nil {
			continue
		}
		messages = append(messages, message)
	}

	return respondJSON(h, 200, messages)
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
