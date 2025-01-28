package main

import (
	"fmt"
	"go-web-bastiario/backend/game"
	"go-web-bastiario/backend/rdb"
	"log"

	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(request *http.Request) bool { return true },
}

func handleConnection(writer http.ResponseWriter, reader *http.Request) {
	connectionHandler, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}
	defer connectionHandler.Close()

	for client := range rdb.ClientSockets {
		log.Println(client)
		x, y, err := rdb.GetPlayerPosition(client)
		if err != nil {
			log.Println("Get error: ", err)
			continue
		}
		msg := []byte(fmt.Sprintf(`{"type": "pl_init", "id": "%s", "x": %d, "y": %d}`, client, x, y))
		if err := connectionHandler.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error: ", err)
			break
		}
	}

	playerObj := game.InitPlayer()

	rdb.AddPlayerSocket(playerObj.ID, connectionHandler)
	rdb.InitPlayerPosition(playerObj.ID, playerObj.X, playerObj.Y)

	defer rdb.RemovePlayer(playerObj.ID)
	defer rdb.DelPlayerSocket(playerObj.ID)

	for {
		var msg []byte
		if _, msg, err = connectionHandler.ReadMessage(); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Printf("Player %v disconnected\n", playerObj.ID)
				return
			} else {
				log.Printf("Message read error from player %v: %v", playerObj.ID, err)
				return
			}
		} else {
			log.Printf("Message received: %s\n", msg)
		}
	}
}

func main() {
	rdb.InitRedis()

	go rdb.SubscribeToChannel("game_updates")

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnection)

	http.Handle("./frontend/assets", http.StripPrefix("./frontend/assets", http.FileServer(http.Dir("./frontend/assets"))))

	port := ":8080"
	fmt.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
