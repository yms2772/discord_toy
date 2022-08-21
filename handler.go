package main

import (
	"fmt"
	"strings"
	"time"
	"toy/api/localtime"
	"toy/api/logger"
	"toy/api/websocket"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	commands := strings.Split(m.Content, " ")

	if len(commands) == 0 {
		return
	}

	logger.Info(m.Content)

	switch commands[0] {
	case "~dday":
		if len(commands) < 4 {
			return
		}

		switch commands[1] {
		case "add":
			_, err := time.Parse("20060102", commands[2])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "날짜 형식은 20220101 입니다. (YYYYMMDD)")

				return
			}

			ddayContent := strings.Join(commands[3:], " ")

			target, err := websocket.AddDDayInfo(m.GuildID, m.Message.Author.ID, commands[2], ddayContent)
			if err != nil {
				// ERROR

				return
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("https://toy.mokky.kr/dday/%s/%s", m.GuildID, target))
		}
	case "~time":
		s.ChannelMessageSend(m.ChannelID, localtime.NowTime().String())
	}
}
