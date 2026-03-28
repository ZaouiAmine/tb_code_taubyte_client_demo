package lib

import (
	"encoding/json"
	"fmt"
	"time"

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

//export handleChatEvents
func handleChatEvents(e event.Event) uint32 {
	pubsubEvent, err := e.PubSub()
	if err != nil {
		return 1
	}

	channel, err := pubsubEvent.Channel()
	if err != nil {
		return 1
	}
	if channel != "chat" {
		return 0
	}

	data, err := pubsubEvent.Data()
	if err != nil {
		return 1
	}

	var message ChatMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return 1
	}
	if message.RoomID == "" || message.Content == "" {
		return 1
	}

	if message.ID == "" {
		message.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}
	if message.CreatedAt == "" {
		message.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	db, err := database.New(dbMatch)
	if err != nil {
		return 1
	}

	encoded, err := json.Marshal(message)
	if err != nil {
		return 1
	}

	key := "chatmsg/" + message.RoomID + "/" + message.ID
	if err := db.Put(key, encoded); err != nil {
		return 1
	}

	return 0
}
