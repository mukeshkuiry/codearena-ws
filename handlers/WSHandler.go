package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/CodeArena-Org/codearena-ws/db"
	"github.com/CodeArena-Org/codearena-ws/helpers"
	"github.com/CodeArena-Org/codearena-ws/models"
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WSHandler_TEMP")
	fmt.Println("DB Rooms: ", db.Rooms)
	fmt.Println("DB FullRooms: ", db.FullRooms)

	// Request Origin API Validation
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey != "123" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API Key: " + apiKey))
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := models.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Read user data from WebSocket connection
	var UserData models.UserModel
	err = conn.ReadJSON(&UserData)
	if err != nil {
		conn.WriteMessage(1, []byte("Invalid user data"))
		log.Println("Failed to read user data:", err)
		return
	}

	// Validate user data
	err = helpers.ValidateUserData(UserData)
	if err != nil {
		conn.WriteMessage(1, []byte("Invalid user data"))
		log.Println("Invalid user data:", err)
		return
	}

	// Check if user is already in queue or room
	if isUserInQueue(UserData.ID) {
		conn.WriteMessage(1, []byte("Already in queue"))
		log.Println("User is already in queue")
		return
	}

	// Find matching room or create a new one
	room, err := helpers.FindRoom(UserData)
	if err != nil {
		newRoom := helpers.CreateRoom(UserData, conn)
		db.Rooms = append(db.Rooms, newRoom)
	} else {
		room.Users = append(room.Users, UserData)
		room.UserConn[conn] = UserData.ID
		db.FullRooms = append(db.FullRooms, *room)
		removeRoomFromDB(room.RoomId)
		helpers.StartGame(*room)
	}

	// Handle messages from the client
	handleMessages(conn)
}

func isUserInQueue(userID string) bool {
	for _, room := range db.Rooms {
		for _, user := range room.Users {
			if user.ID == userID {
				return true
			}
		}
	}

	return false
}

func removeRoomFromDB(roomID string) {
	for i, r := range db.Rooms {
		if r.RoomId == roomID {
			db.Rooms = append(db.Rooms[:i], db.Rooms[i+1:]...)
			break
		}
	}
}

func handleMessages(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			return
		}
		log.Printf("Message Received: %s, Message Type: %d\n", msg, msgType)
	}
}
