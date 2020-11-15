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

	if err := resolver.RunDistributionStrategy(); err != nil {
		fmt.Printf("Error while running distribution strategy: %v\n", err)
	}

	//// Pick a random FPL team for analysis
	//team, err := resolver.PickRandomTeam()
	//if err != nil {
	//	fmt.Printf("Error occured while creating a team: %v\n", err)
	//	return
	//}
	//fmt.Println(team)
	//
	//costDistribution := strategy.CalculateTeamDistribution(team)
	//
	//fmt.Printf("Level One Players: %v\n", costDistribution[0])
	//fmt.Printf("Level Two Players: %v\n", costDistribution[1])
	//fmt.Printf("Level Three Players: %v\n", costDistribution[2])
	//
	//points, err := resolver.CalculateTeamPoints(team)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Printf("Total Points: %v\n", points)
}
