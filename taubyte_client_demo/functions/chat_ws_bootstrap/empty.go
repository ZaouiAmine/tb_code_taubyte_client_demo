package lib

import (
	"encoding/json"

	"github.com/taubyte/go-sdk/event"
	httpevent "github.com/taubyte/go-sdk/http/event"
	pubsubnode "github.com/taubyte/go-sdk/pubsub/node"
)

const chatChannel = "chat"

type websocketBootstrapResponse struct {
	Channel      string `json:"channel"`
	WebSocketURL string `json:"websocketUrl"`
}

//export getChatWebsocket
func getChatWebsocket(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	channel, err := pubsubnode.Channel(chatChannel)
	if err != nil {
		return respondError(h, 500, "failed to get pubsub channel")
	}

	url, err := channel.WebSocket().Url()
	if err != nil {
		return respondError(h, 500, "failed to generate websocket url")
	}

	resp := websocketBootstrapResponse{
		Channel:      chatChannel,
		WebSocketURL: url.String(),
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

func respondError(h httpevent.Event, status int, message string) uint32 {
	type errorResponse struct {
		Error string `json:"error"`
	}
	return respondJSON(h, status, errorResponse{Error: message})
}
