package lib

import (
	"encoding/json"

	"github.com/taubyte/go-sdk/event"
	httpevent "github.com/taubyte/go-sdk/http/event"
	pubsubnode "github.com/taubyte/go-sdk/pubsub/node"
)

const chatChannel = "chat"

type websocketBootstrapResponse struct {
	Channel           string `json:"channel"`
	WebSocketURL      string `json:"websocketUrl"`
	RealtimeAvailable bool   `json:"realtimeAvailable"`
	Message           string `json:"message,omitempty"`
}

//export getChatWebsocket
func getChatWebsocket(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	resp := websocketBootstrapResponse{
		Channel:           chatChannel,
		RealtimeAvailable: false,
		Message:           "realtime unavailable, use polling",
	}

	channel, err := pubsubnode.Channel(chatChannel)
	if err == nil {
		url, urlErr := channel.WebSocket().Url()
		if urlErr == nil {
			resp.WebSocketURL = url.String()
			resp.RealtimeAvailable = true
			resp.Message = ""
		}
	}

	return respondJSON(h, 200, resp)
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
