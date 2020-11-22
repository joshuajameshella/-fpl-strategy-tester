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

	// Simulate the teams used to feed into the different FPL strategies
	simulatedTeams, errs := resolver.GenerateTeams()
	for err := range errs {
		fmt.Printf("Error while simulating teams: %v\n", err)
	}

	// Run the cost variation strategy
	if err := resolver.RunCostVariationStrategy(simulatedTeams); err != nil {
		fmt.Printf("Error while running distribution strategy: %v\n", err)
	}

	// Run the distribution strategy
	//if err := resolver.RunDistributionStrategy(); err != nil {
	//	fmt.Printf("Error while running distribution strategy: %v\n", err)
	//}
}
