package db

import "github.com/CodeArena-Org/codearena-ws/models"

var Rooms = make([]models.Room, 0)

var FullRooms = make([]models.Room, 0)

// types of events during any game

// start the game, someone left, acceted, fail, 