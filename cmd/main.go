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

	// Run the distribution strategy
	//if err := resolver.RunDistributionStrategy(); err != nil {
	//	fmt.Printf("Error while running distribution strategy: %v\n", err)
	//}

	if err := resolver.RunCostVariationStrategy(); err != nil {
		fmt.Printf("Error while running distribution strategy: %v\n", err)
	}

}
