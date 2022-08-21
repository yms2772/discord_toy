package apps

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"nhooyr.io/websocket"
)

type DDay struct {
	app.Compo

	upgrader *websocket.Conn

	guildID     string
	path        string
	date        string
	remain      string
	description string
}

type ddayJSON struct {
	Date        string `json:"date"`
	Description string `json:"description"`
}

type receive struct {
	ddayJSON

	Type    string `json:"type"`
	Content string `json:"content"`
}
