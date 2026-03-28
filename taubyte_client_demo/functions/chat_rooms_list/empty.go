package lib

import (
	"encoding/json"
	"fmt"
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
		return writeError(h, 500, "failed to connect database")
	}

	keys, err := db.List("chatroom/")
	if err != nil {
		return writeError(h, 500, "failed to list chat rooms")
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

	payload, err := json.Marshal(rooms)
	if err != nil {
		return writeError(h, 500, "failed to encode rooms")
	}

	h.Headers().Set("Content-Type", "application/json")
	h.Write(payload)
	h.Return(200)
	return 0
}

func writeError(h interface {
	Write([]byte) (int, error)
	Return(int) error
}, status int, message string) uint32 {
	h.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", message)))
	h.Return(status)
	return 0
}
