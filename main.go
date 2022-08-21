package main

import (
	"fmt"
	"net/http"
	"strings"
	"toy/api/logger"
	"toy/api/src"
	"toy/api/websocket"
	"toy/apps"

	"github.com/bwmarrin/discordgo"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")

	if len(path) > 2 {
		switch path[1] {
		case "dday":
			m.DDay.ServeHTTP(w, r)
		case "api":
			m.API.ServeHTTP(w, r)
		default:
			m.DDay.ServeHTTP(w, r)
		}
	} else {
		m.DDay.ServeHTTP(w, r)
	}
}

func main() {
	app.RouteWithRegexp("/dday/[0-9]+/[a-zA-Z0-9_-]+", &apps.DDay{})
	app.RunWhenOnBrowser()

	port := src.GetValue("webserver", "port")
	discordToken = src.GetValue("token", "discord_token")

	if len(discordToken) == 0 || discordToken == "none" {
		logger.Error("'discord_token' in 'setting.ini' does not exist")

		return
	}

	logger.Info("configuration file has been loaded")

	logger.Info("authenticating discord api...")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		logger.Error(err.Error())

		return
	}
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	discord.AddHandler(messageCreate)

	if err := discord.Open(); err != nil {
		logger.Error(err.Error())

		return
	}

	mux := &Mux{
		DDay: http.NewServeMux(),
		API:  http.NewServeMux(),
	}

	mux.DDay.Handle("/", &app.Handler{
		Title: "D-Day",
		Styles: []string{
			"/web/dday.css",
		},
		Icon: app.Icon{
			Default:    "/web/favicon/192.png",
			Large:      "/web/favicon/512.png",
			AppleTouch: "/web/favicon/192.png",
		},
		LoadingLabel: "getting dday info...",
	})
	mux.API.HandleFunc("/api/", websocket.Websocket)

	logger.Info(fmt.Sprintf("webserver is opened on '%s' port", port))
	_ = http.ListenAndServe(":"+port, mux)
}
