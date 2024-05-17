package helpers

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/CodeArena-Org/codearena-ws/db"
	"github.com/CodeArena-Org/codearena-ws/models"
)

func FindRoom(UserData models.UserModel) (*models.Room, error) {

	// find matching room with higest match score
	var maxMatchScore float64 = 0
	var maxMatchRoom *models.Room = nil

	for _, room := range db.Rooms {
		if len(room.Users) == 1 {
			// check if room is suitable for his level and timeslot
			if err := ValidateRoomType(room, UserData); err != nil {
				continue
			}

			matchScore := calculateMatchScore(UserData, room.Users[0])
			if matchScore > maxMatchScore {
				maxMatchScore = matchScore
				maxMatchRoom = &room
			}
		}
	}

	if maxMatchRoom == nil {
		return nil, errors.New("no match found")
	}

	return maxMatchRoom, nil
}

func ValidateRoomType(room models.Room, UserData models.UserModel) error {
	// check if room is suitable for his level and timeslot

	if UserData.ProblemLevel != room.Users[0].ProblemLevel {
		return errors.New("problem level mismatch")
	}
	if UserData.TimeSlot != room.Users[0].TimeSlot {
		return errors.New("time slot mismatch")
	}

	if math.Abs(UserData.Rating-room.Users[0].Rating) > 100 {
		return errors.New("rating mismatch")
	}
	return nil
}

func calculateMatchScore(user1 models.UserModel, user2 models.UserModel) float64 {
	matchScore := 0.0

	ratingDiff := math.Abs(float64(user1.Rating - user2.Rating))
	dsaPreferenceScore := calculateDSAPreferenceMatch(user1, user2)

	matchScore = 10 - math.Min(10, ratingDiff/10) + dsaPreferenceScore

	return matchScore
}

func calculateDSAPreferenceMatch(user1 models.UserModel, user2 models.UserModel) float64 {
	matchScore := 0.0
	for _, dsa := range user1.DSAPreferences {
		for _, dsa2 := range user2.DSAPreferences {
			if dsa == dsa2 {
				matchScore += 1
			}
		}
	}
	return matchScore
}

func CreateRoom(UserData models.UserModel, conn *websocket.Conn) models.Room {
	room := models.Room{
		RoomId: createRandomRoomID(),
		Users:  []models.UserModel{UserData},
		UserConn: map[*websocket.Conn]string{
			conn: UserData.ID,
		},
	}
	return room
}

func createRandomRoomID() string {
	roomID := "room"
	roomID += strconv.Itoa(rand.Intn(100000))
	return roomID
}

func ValidateUserData(UserData models.UserModel) error {
	if UserData.Name == "" {
		return fmt.Errorf("invalid user name")
	}
	if UserData.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if UserData.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if UserData.Rating < 0 {
		return fmt.Errorf("invalid rating")
	}
	if UserData.ProblemLevel != models.EasyProblemLevel && UserData.ProblemLevel != models.MediumProblemLevel && UserData.ProblemLevel != models.HardProblemLevel {
		return fmt.Errorf("invalid problem level")
	}
	if UserData.TimeSlot != models.TenMinuteSlot && UserData.TimeSlot != models.ThirtyMinuteSlot && UserData.TimeSlot != models.SixtyMinuteSlot {
		return fmt.Errorf("invalid time slot")
	}
	return nil
}

func StartGame(room models.Room) {
	// send message to all users in room
	for conn := range room.UserConn {
		conn.WriteMessage(websocket.TextMessage, []byte("Game started"))
	}
}
