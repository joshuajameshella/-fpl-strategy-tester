package strategy

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math"
	"sort"
)

/*	DISTRIBUTION:
	This file of code manages the simulation of team cost distribution.
 	The theory is that higher priced players will yield more points, and that
	a ideal cost distribution of cheap to expensive players exists.
*/

// CalculateTeamDistribution takes the team and calculates what tier each player is
func CalculateTeamDistribution(team []database.PlayerInfo) []int {

	// Level 1, Level 2, Level 3
	costDistribution := []int{0, 0, 0}

	// Calculate cost distribution of each player
	for _, goalkeeper := range team[0:2] {
		if goalkeeper.Price == 60 {
			costDistribution[0]++
			continue
		}
		if goalkeeper.Price > 45 {
			costDistribution[1]++
			continue
		}
		costDistribution[2]++
	}
	for _, defender := range team[2:7] {
		if defender.Price >= 65 {
			costDistribution[0]++
			continue
		}
		if defender.Price > 50 {
			costDistribution[1]++
			continue
		}
		costDistribution[2]++
	}
	for _, midfielder := range team[7:12] {
		if midfielder.Price >= 90 {
			costDistribution[0]++
			continue
		}
		if midfielder.Price > 65 {
			costDistribution[1]++
			continue
		}
		costDistribution[2]++
	}

	for _, forward := range team[12:15] {
		if forward.Price >= 90 {
			costDistribution[0]++
			continue
		}
		if forward.Price > 65 {
			costDistribution[1]++
			continue
		}
		costDistribution[2]++
	}

	return costDistribution
}

// ProcessDistributionResults uses the simulation data to create values used in plotting charts
func ProcessDistributionResults(resultsCh chan []int) {

	// Create an array to house each category of distribution, between 0 and 10
	distributionResults := make([][]int, 10)

	// For each result simulated, store result in the correct array space
	for result := range resultsCh {
		distributionResults[result[0]] = append(distributionResults[result[0]], result[1])
	}

	// Print results to user
	for key, category := range distributionResults {
		percentiles := findPercentiles(category)
		fmt.Printf("Distribution for %v valuable player: %v\n", key, percentiles)
	}
}

// findPercentiles takes the points array for each team distribution and returns the
// necessary percentiles. These will be used to plot the box charts.
func findPercentiles(distributionData []int) []int {

	// If there isn't enough data to display a full box plot, ignore
	if len(distributionData) < 5 {
		return nil
	}

	sort.Ints(distributionData)
	percentiles := []int{
		distributionData[int(math.Floor(float64(len(distributionData))*0.05))],
		distributionData[int(math.Floor(float64(len(distributionData))*0.25))],
		distributionData[int(math.Floor(float64(len(distributionData))*0.50))],
		distributionData[int(math.Floor(float64(len(distributionData))*0.75))],
		distributionData[int(math.Floor(float64(len(distributionData))*0.95))],
	}

	return percentiles
}
