package models

import "github.com/gorilla/websocket"

const EasyProblemLevel = "easy"
const MediumProblemLevel = "medium"
const HardProblemLevel = "hard"

const TenMinuteSlot = 10
const ThirtyMinuteSlot = 30
const SixtyMinuteSlot = 60

type UserModel struct {
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	ID             string   `json:"id"`
	Rating         float64  `json:"rating"`
	DSAPreferences []string `json:"dsaPreferences"`
	ProblemLevel   string   `json:"problemLevel"`
	TimeSlot       int      `json:"timeSlot"`
}

type Room struct {
	RoomId    string
	Users     []UserModel
	UserConn map[*websocket.Conn]string
}
