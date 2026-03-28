package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
)

type ChatRoom struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

const dbMatch = "appdata"

//export createChatRoom
func createChatRoom(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	body, err := io.ReadAll(h.Body())
	if err != nil {
		return respondError(h, 400, "failed to read request body")
	}
	defer h.Body().Close()

	var room ChatRoom
	if err := json.Unmarshal(body, &room); err != nil {
		return respondError(h, 400, "invalid chat room payload")
	}
	if room.Name == "" {
		return respondError(h, 400, "room name is required")
	}

	room.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	room.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	encoded, err := json.Marshal(room)
	if err != nil {
		return respondError(h, 500, "failed to encode chat room")
	}
	if err := db.Put("chatroom/"+room.ID, encoded); err != nil {
		return respondError(h, 500, "failed to store chat room")
	}

	return respondJSON(h, 201, room)
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
