package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math"
	"sort"
	"sync"

	"github.com/icelolly/go-errors"
)

/*	DISTRIBUTION:
	This file of code manages the simulation of team cost distribution.
 	The theory is that higher priced players will yield more points, and that
	a ideal cost distribution of cheap to expensive players exists.
*/

// RunDistributionStrategy simulates random teams, and records the points and cost distribution for each
func (r *Resolver) RunDistributionStrategy() error {

	// How many simulations to run
	const maxQueries int = 100
	const maxBatchSize int = 100

	// Data channels used to store simulation simulation_results
	resultsCh := make(chan []int, maxQueries)
	errCh := make(chan error, maxQueries)

	// Empty the results file of any old data, and write the headers to the file
	resultsFilePath := "internal/simulation_results/cost_distribution.csv"
	if err := truncateFile(resultsFilePath); err != nil {
		fmt.Println(err)
	}
	if err := writeToFile(
		resultsFilePath,
		fmt.Sprintf("Tier One Players, 5th Percentile, 25th Percentile, 50th Percentile, 75th Percentile, 95th Percentile\n")); err != nil {
		return errors.Wrap(err)
	}

	// Simulate distribution strategy in batches (Prevent MySQL connection error 1040)
	for j := 0; j < (maxQueries / maxBatchSize); j++ {

		// Manage concurrency
		wg := &sync.WaitGroup{}
		wg.Add(maxBatchSize)

		for i := 0; i < maxBatchSize; i++ {
			go func() {
				defer wg.Done()

				// Simulate a random FPL team
				team, err := r.PickRandomTeam(950)
				if err != nil {
					errCh <- err
				}

				// Calculate the overall team points
				points, err := r.CalculateTeamPoints(team)
				if err != nil {
					errCh <- err
				}

				// Calculate the cost distribution of the team, then add data into appropriate channels
				if costDistribution, err := CalculateTeamDistribution(team); err == nil {
					resultsCh <- []int{costDistribution[0], points}
				} else {
					errCh <- err
				}

			}()
		}
		wg.Wait()
	}

	close(errCh)
	close(resultsCh)

	// Log any errors which may have occurred
	for err := range errCh {
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	// Handle the simulation simulation_results
	// percentiles := ProcessDistributionResults(resultsCh)

	return nil
}

// RunCostVariationStrategy takes the simulated random teams over a range of total values in order to determine
// the relationship between cost and points
func (r *Resolver) RunCostVariationStrategy(simulatedTeams chan []database.PlayerInfo) error {

	//// Empty the results file of any old data, and write the headers to the file
	//resultsFilePath := "internal/simulation_results/cost_distribution.csv"
	//if err := truncateFile(resultsFilePath); err != nil {
	//	return errors.Wrap(err)
	//}
	//if err := writeToFile(resultsFilePath, fmt.Sprintf("Team Price, Team Points\n")); err != nil {
	//	return errors.Wrap(err)
	//}

	// Process each team
	dataStr := make([]string, 0)
	for team := range simulatedTeams {

		// Calculate the overall team points
		points, err := r.CalculateTeamPoints(team)
		if err != nil {
			fmt.Println(err)
		}

		// Calculate the overall team price
		teamPrice := CalculatePrice(team)

		dataStr = append(dataStr, fmt.Sprintf("%v, %v", teamPrice, points))
	}

	//// Write the results to the 'cost_distribution' file
	//if err := writeToFile(resultsFilePath, strings.Join(dataStr, "\n")); err != nil {
	//	return errors.Wrap(err)
	//}

	return nil
}

// CalculateTeamDistribution takes the team and calculates what tier each player is
func CalculateTeamDistribution(team []database.PlayerInfo) ([]int, error) {

	// Level 1, Level 2, Level 3
	costDistribution := []int{0, 0, 0}

	if len(team) < 15 {
		return nil, errors.New("Empty team returned")
	}

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

	return costDistribution, nil
}

// ProcessDistributionResults uses the simulation data to create values used in plotting charts
func ProcessDistributionResults(resultsCh chan []int) [][]int {

	// Create an array to house each category of distribution, between 0 and 10
	distributionResults := make([][]int, 10)

	// For each result simulated, store result in the correct array space
	for result := range resultsCh {
		distributionResults[result[0]] = append(distributionResults[result[0]], result[1])
	}

	// Print simulation_results to user
	percentiles := make([][]int, 10)
	for key, category := range distributionResults {
		percentiles[key] = findPercentiles(category)
		fmt.Printf("Distribution for %v valuable player: %v\n", key, percentiles)
	}

	return percentiles
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
