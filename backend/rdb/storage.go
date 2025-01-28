package rdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	ctx                = context.Background()
	db                 *redis.Client
	clientSockets      map[string]*websocket.Conn
	clientSocketsMutex sync.Mutex
)

func InitRedis() {
	db = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	clientSockets = make(map[string]*websocket.Conn)

	if _, err := db.Ping(ctx).Result(); err != nil {
		log.Fatal("!@#$% Redis connection failed: ", err)
	} else {
		log.Println("Connected to Redis")
	}
}

func AddPlayerSocket(id string, socket *websocket.Conn) {
	clientSocketsMutex.Lock()
	defer clientSocketsMutex.Unlock()

	clientSockets[id] = socket
}

func DelPlayerSocket(id string) {
	clientSocketsMutex.Lock()
	defer clientSocketsMutex.Unlock()

	delete(clientSockets, id)
}

func publishMessage(channel string, message string) {
	err := db.Publish(ctx, channel, message).Err()
	if err != nil {
		log.Println("Publish error: ", err)
	}
}

func SubscribeToChannel(channel string) {
	pubsub := db.Subscribe(ctx, channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("Subscription error: ", err)
			return
		}
		log.Printf("Received message from channel %s: %s", msg.Channel, msg.Payload)

		for client, socket := range clientSockets {
			err := socket.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				socket.Close()
				delete(clientSockets, client)
			}
		}
	}
}

func InitPlayerPosition(playerID string, x, y int) {
	err := db.HMSet(ctx, "player:"+playerID, []string{
		"x", fmt.Sprintf("%d", x),
		"y", fmt.Sprintf("%d", y),
	}).Err()
	if err != nil {
		log.Println("!@#$% Failed to save player position:", err)
	} else {
		publishMessage("game_updates", fmt.Sprintf(`{"type": "pl_init", "id": "%s", "x": %d, "y": %d}`, playerID, x, y))
	}
}

func ChangePlayerPosition(playerID string, x, y int) {
	err := db.HMSet(ctx, "player:"+playerID, []string{
		"x", fmt.Sprintf("%d", x),
		"y", fmt.Sprintf("%d", y),
	}).Err()
	if err != nil {
		log.Println("!@#$% Failed to save player position:", err)
	} else {
		publishMessage("game_updates", fmt.Sprintf(`{"type": "pl_chng_loc", "id": "%s", "x": %d, "y": %d}`, playerID, x, y))
	}
}

func GetPlayerPosition(playerID string) (int, int, error) {
	result, err := db.HGetAll(ctx, "player:"+playerID).Result()
	if err != nil {
		return 0, 0, err
	}

	x, _ := strconv.Atoi(result["x"])
	y, _ := strconv.Atoi(result["y"])
	return x, y, nil
}

func RemovePlayer(playerID string) {
	err := db.Del(ctx, "player:"+playerID).Err()
	if err != nil {
		log.Println("Failed to remove player:", err)
	} else {
		publishMessage("game_updates", fmt.Sprintf(`{"type": "pl_del", "id": "%s"}`, playerID))
	}
}
