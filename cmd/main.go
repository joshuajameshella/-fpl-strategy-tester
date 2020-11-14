package main

import (
	"fmt"
	"fpl-strategy-tester/internal"
	"math/rand"
	"time"
)

func main() {

	resolver := internal.NewResolver()
	resolver.ResolveDatabase()

	// Set the seed used for generating random numbers
	rand.Seed(time.Now().UnixNano())

	// Pick a random FPL team for analysis
	team, err := resolver.PickRandomTeam()
	if err != nil {
		fmt.Printf("Error occured while creating a team: %v\n", err)
		return
	}
	fmt.Println(team)
}
