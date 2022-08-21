package websocket

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"toy/api/localtime"
	"toy/api/logger"
	"toy/api/src"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setUpgrader(w http.ResponseWriter, r *http.Request, guildID string, sessionID int64) (err error) {
	for _, i := range conn[guildID] {
		if i.SessionID == sessionID {
			i.Upgrader, err = upgrader.Upgrade(w, r, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func closeUpgrader(guildID string, sessionID int64) (err error) {
	for _, i := range conn[guildID] {
		if i.SessionID == sessionID {
			if err = i.Upgrader.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func removeUpgrader(guildID string, sessionID int64) {
	for i := 0; i < len(conn); i++ {
		c := conn[guildID][i]
		if c.SessionID == sessionID {
			conn[guildID] = append(conn[guildID][:i], conn[guildID][i+1:]...)
			i--

			continue
		}
	}
}

func sendWebsocket(messageType int, guildID string, sessionID int64, data map[string]interface{}) {
	sendData, _ := json.Marshal(data)

	for _, i := range conn[guildID] {
		if i == nil {
			continue
		}

		if sessionID == 0 || i.SessionID == sessionID {
			if err := i.Upgrader.WriteMessage(messageType, sendData); err != nil {
				continue
			}
		}
	}
}

func connDB() (client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(src.GetValue("database", "mongodb_url")).SetAuth(options.Credential{
		Username: src.GetValue("database", "mongodb_id"),
		Password: src.GetValue("database", "mongodb_pw"),
	}))

	return
}

func getCollection(client *mongo.Client, db, collection string) *mongo.Collection {
	return client.Database(db).Collection(collection)
}

func AddDDayInfo(guildID, creator, date, content string) (target string, err error) {
	client, ctx, cancel := connDB()
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}(client, ctx)
	defer cancel()

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d|%s|%s|%s", localtime.NowTime().UnixNano(), creator, date, content)))

	target = hex.EncodeToString(hash.Sum(nil))

	getCollection(client, "dday", "ddaylist").InsertOne(ctx, bson.M{
		"target":      target,
		"guild_id":    guildID,
		"creator":     creator,
		"date":        date,
		"description": content,
	})

	return target, nil
}

func getDDayInfo(guildID string, target string) (data *dday, err error) {
	client, ctx, cancel := connDB()
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}(client, ctx)
	defer cancel()

	var d dday

	if err := getCollection(client, "dday", "ddaylist").FindOne(ctx, bson.M{
		"target":   target,
		"guild_id": guildID,
	}).Decode(&d); err != nil {
		logger.Error(err.Error())
	}

	return &d, err
}
