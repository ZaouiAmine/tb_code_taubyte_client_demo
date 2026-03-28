package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	httpevent "github.com/taubyte/go-sdk/http/event"
	pubsubnode "github.com/taubyte/go-sdk/pubsub/node"
)

type ChatMessage struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

const (
	dbMatch     = "appdata"
	chatChannel = "chat"
)

//export postChatMessage
func postChatMessage(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	body, err := io.ReadAll(h.Body())
	if err != nil {
		return respondError(h, 400, "failed to read request body")
	}
	defer h.Body().Close()

	var message ChatMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return respondError(h, 400, "invalid chat payload")
	}
	if message.RoomID == "" || message.Content == "" {
		return respondError(h, 400, "roomId and content are required")
	}

	now := time.Now().UTC()
	message.ID = fmt.Sprintf("%d", now.UnixNano())
	message.CreatedAt = now.Format(time.RFC3339)

	db, err := database.New(dbMatch)
	if err != nil {
		return respondError(h, 500, "failed to connect database")
	}

	key := "chatmsg/" + message.RoomID + "/" + message.ID
	encoded, err := json.Marshal(message)
	if err != nil {
		return respondError(h, 500, "failed to encode chat message")
	}
	if err := db.Put(key, encoded); err != nil {
		return respondError(h, 500, "failed to store chat message")
	}

	channel, err := pubsubnode.Channel(chatChannel)
	if err == nil {
		_ = channel.Publish(encoded)
	}

	return respondJSON(h, 201, message)
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
