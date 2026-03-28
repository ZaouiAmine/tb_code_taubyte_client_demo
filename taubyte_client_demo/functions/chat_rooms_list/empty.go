package lib

import (
	"encoding/json"
	"sort"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
)

type ChatRoom struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

const dbMatch = "appdata"

//export listChatRooms
func listChatRooms(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	keys, err := db.List("chatroom/")
	if err != nil {
		return respondError(h, 500, "failed to list chat rooms")
	}
	sort.Strings(keys)

	rooms := make([]ChatRoom, 0, len(keys))
	for _, key := range keys {
		raw, getErr := db.Get(key)
		if getErr != nil {
			continue
		}

		var room ChatRoom
		if unmarshalErr := json.Unmarshal(raw, &room); unmarshalErr != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	return respondJSON(h, 200, rooms)
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
