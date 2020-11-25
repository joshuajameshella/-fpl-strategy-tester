package main

import (
	"fpl-strategy-tester/internal"
	"fpl-strategy-tester/internal/database"
	"log"
	"math/rand"
	"time"
)

func main() {

	resolver := internal.NewResolver()
	resolver.ResolveDatabase()
	resolver.ResolveCache()

	// Set the seed used for generating random numbers
	rand.Seed(time.Now().UnixNano())

	// Data channels used to store simulation simulation_results
	resultsCh := make(chan []database.PlayerInfo, internal.MaxQueries)
	errCh := make(chan error, internal.MaxQueries)

	// Simulate the teams used to feed into the different FPL strategies
	log.Printf("-> Simulating 10,000 random FPL teams...\t")
	resolver.GenerateTeams(resultsCh, errCh)

	// Run the cost variation strategy
	log.Printf("-> Running Cost Variation strategy...\t")
	if err := resolver.RunCostVariationStrategy(resultsCh); err != nil {
		log.Printf("Error: %v\n", err)
	}

	//Run the cost distribution strategy
	log.Printf("-> Running Cost Distribution strategy...\t")
	if err := resolver.RunDistributionStrategy(resultsCh); err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Close the channels and process any errors
	defer close(resultsCh)
	close(errCh)
	for err := range errCh {
		log.Printf("Error: %v\n", err)
	}
}
