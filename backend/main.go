package main

import (
	"fmt"
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

	log.Println("New connection")

	for {
		var msg []byte
		if _, msg, err = connectionHandler.ReadMessage(); err != nil {
			log.Println("!@#$% Message read error: ", err)
			return
		}

		log.Printf("Message received: %s\n", msg)

		if err := connectionHandler.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("!@#$% Message write error: ", err)
			return
		}
	}
}

func main() {
	initRedis()

	saveKeyVal("42", "28")

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnection)

	http.Handle("./frontend/assets", http.StripPrefix("./frontend/assets", http.FileServer(http.Dir("./frontend/assets"))))

	port := ":8080"
	fmt.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
