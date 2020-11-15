package database

import (
	"github.com/doug-martin/goqu/v9"
)

// dataGW1 is the database table used to store the pre-season player data
var dataGW1 = goqu.T("GW1")

// playerData is the database table used to store the player data by game week
var playerData = goqu.T("GW_data")

// Database login info
const dbSchemaName = "fpl"
const dbAddress = "127.0.0.1:3306"
const dbUsername = "root"

// Credentials is the data structure found in 'credentials.json'
type Credentials struct {
	DBUsername string
	DBPassword string
}

// PlayerInfo is the structure of data found in the 'GW1' table
type PlayerInfo struct {
	ID        int
	FirstName string
	LastName  string
	Position  string
	Price     int
	Team      int
}

// PlayerGWInfo is the structure of data found in the 'GW_data' table
type PlayerGWInfo struct {
	Name         string
	Element      int
	OpponentTeam int
	TotalPoints  int
	Value        int
	WasHome      string
	GW           int
}

// Constant time format to be used throughout project
const TimeFormat = "2006-01-02 15:04:05"
