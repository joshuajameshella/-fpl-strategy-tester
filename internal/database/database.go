package database

import (
	"github.com/doug-martin/goqu/v9"
)

// dataGW1 is the database table used to store the pre-season player data
var dataGW1 = goqu.T("GW1")

// Database login info
const dbSchemaName = "fpl"
const dbAddress = "127.0.0.1:3306"
const dbUsername = "root"

// Constant time format to be used throughout project
const TimeFormat = "2006-01-02 15:04:05"
