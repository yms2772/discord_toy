package websocket

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type receive struct {
	Type string `json:"type"`
}

type dday struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Target      string             `json:"target" bson:"target"`
	GuildID     string             `json:"guild_id" bson:"guild_id"`
	Creator     string             `json:"creator" bson:"creator"`
	Date        string             `json:"date" bson:"date"`
	Description string             `json:"description" bson:"description"`
}

type websocketConn struct {
	Upgrader  *websocket.Conn
	SessionID int64
}
