package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"toy/api/logger"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"nhooyr.io/websocket"
)

func (d *DDay) SendWebsocket(messageType websocket.MessageType, data map[string]interface{}) error {
	if d.upgrader == nil {
		return errors.New("connection does not exist")
	}

	ctx := context.Background()
	sendData, _ := json.Marshal(data)

	err := d.upgrader.Write(ctx, messageType, sendData)
	if err != nil {
		return err
	}

	return nil
}

func (d *DDay) OnMount(ctx app.Context) {
	var err error

	domain := app.Window().URL().Host
	paths := strings.Split(app.Window().URL().Path, "/")

	if len(paths) < 4 {
		log.Println(paths)

		return
	}

	d.guildID = paths[2]
	d.path = paths[3]

	if len(d.guildID) == 0 || len(d.path) == 0 {
		return
	}

	// TODO: domain change
	d.upgrader, _, err = websocket.Dial(ctx, "wss://"+domain+"/api/"+d.guildID+"/"+d.path, nil)
	if err != nil {
		log.Println(err)

		return
	}

	go func() {
		defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
			err := conn.Close(code, reason)
			if err != nil {
				return
			}
		}(d.upgrader, websocket.StatusInternalError, "StatusInternalError")

		for {
			_, message, err := d.upgrader.Read(ctx)
			if err != nil {
				log.Println(err)
			}

			var r receive
			_ = json.Unmarshal(message, &r)

			switch r.Type {
			case "update_dday":
				loc, err := time.LoadLocation(os.Getenv("TZ"))
				if err != nil {
					loc, _ = time.LoadLocation("Asia/Seoul")
				}

				expiredDate, err := time.ParseInLocation("20060102150405", r.Date+"000000", loc)
				if err != nil {
					logger.Error(err.Error())

					continue
				}

				expiredDateSeconds := int(time.Until(expiredDate).Seconds())

				if expiredDateSeconds < 0 {
					expiredDateSeconds *= -1
				}

				expiredDateH := expiredDateSeconds / 3600
				expiredDateM := (expiredDateSeconds - (3600 * expiredDateH)) / 60
				expiredDateS := expiredDateSeconds - (3600 * expiredDateH) - (expiredDateM * 60)

				remainDay := math.Ceil(time.Since(expiredDate).Hours() / -24)
				remainDayStr := fmt.Sprintf("D-%d", int(remainDay))

				if remainDay < 0 {
					remainDayStr = fmt.Sprintf("D+%d", int(-remainDay))
				}

				d.date = fmt.Sprintf("%s / %s", remainDayStr, expiredDate.Format("2006년 01월 02일"))
				d.remain = fmt.Sprintf("%d시간 %02d분 %02d초", expiredDateH, expiredDateM, expiredDateS)
				d.description = r.Description
				d.Update()

				logger.Info(fmt.Sprintf("set: %s / %s / %s", d.date, d.remain, d.description))
			}
		}
	}()
}

func (d *DDay) OnNav(_ app.Context) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)

		for range ticker.C {
			_ = d.SendWebsocket(websocket.MessageText, map[string]interface{}{
				"type": "update_dday",
			})
		}
	}()
}

func (d *DDay) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			app.Text(d.date),
		).
			Class("date"),
		app.P().Body(
			app.Text(d.remain),
		).
			Class("time"),
		app.P().Body(
			app.Text(d.description),
		).
			Class("text"),
	).
		ID("clock")
}
