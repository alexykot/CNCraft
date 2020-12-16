package status

import (
	"github.com/alexykot/cncraft/pkg/chat"
)

const (
	ServerName = "CNCraft"
	ServerUUID = "41d1fed5-aa44-432c-ab1b-2810001f3270" // TODO supply from config or random generate?

	ServerMotd = "LOLWUT?"
	ServerIcon = ""
)

type Response struct {
	Version     Version `json:"version,string"`
	Players     Players `json:"players,string"`
	Description Message `json:"description"`
	Favicon     string  `json:"favicon"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int            `json:"max"`
	Online int            `json:"online"`
	Sample []SamplePlayer `json:"sample"`
}

type SamplePlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}

func DefaultResponse(currentVersion int) Response {
	return Response{
		Version: Version{
			Name:     "CNCraft Server",
			Protocol: currentVersion,
		},
		Players: Players{
			Max:    10,
			Online: 1,
			Sample: []SamplePlayer{
				{
					Name: ServerName,
					ID:   ServerUUID,
				},
			},
		},
		Description: Message{
			Text: chat.Translate(ServerMotd),
		},
		Favicon: ServerIcon,
	}
}
