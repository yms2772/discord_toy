package websocket

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"toy/api/localtime"
	"toy/api/logger"

	"github.com/gorilla/websocket"
)

func Websocket(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")

	if len(paths) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))

		return
	}

	_, err := strconv.ParseInt(paths[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))

		return
	}

	guildID := paths[2]
	target := paths[3]

	if len(target) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))

		return
	}

	keepAliveTicker := time.NewTicker(10 * time.Second)
	sessionID := localtime.NowTime().UnixNano()

	conn[guildID] = append(conn[guildID], &websocketConn{
		Upgrader:  nil,
		SessionID: sessionID,
	})

	if err := setUpgrader(w, r, guildID, sessionID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))

		return
	}

	defer func() {
		logger.Info("websocket closed")
		_ = closeUpgrader(guildID, sessionID)
		removeUpgrader(guildID, sessionID)
		keepAliveTicker.Stop()
	}()

	go func() {
		for range keepAliveTicker.C {
			sendWebsocket(websocket.TextMessage, guildID, sessionID, map[string]interface{}{
				"type": "ping",
			})
		}
	}()

LOOP:
	for _, i := range conn[guildID] {
		if i.SessionID == sessionID {
			for {
				_, message, err := i.Upgrader.ReadMessage()
				if err != nil {
					break LOOP
				}

				var r receive
				_ = json.Unmarshal(message, &r)

				switch r.Type {
				case "update_dday":
					d, err := getDDayInfo(guildID, target)
					if err != nil {
						logger.Error(err.Error())

						continue
					}

					sendWebsocket(websocket.TextMessage, guildID, sessionID, map[string]interface{}{
						"type":        "update_dday",
						"date":        d.Date,
						"description": d.Description,
					})
				}
			}
		}
	}
}
