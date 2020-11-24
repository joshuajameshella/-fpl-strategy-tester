package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/icelolly/go-errors"
)

/*	DISTRIBUTION:
	This file of code manages the simulation of team cost distribution.
 	The theory is that higher priced players will yield more points, and that
	a ideal cost distribution of cheap to expensive players exists.
*/

// RunCostVariationStrategy takes the simulated random teams over a range of total values in order to determine
// the relationship between cost and points
func (r *Resolver) RunCostVariationStrategy(simulatedTeams chan []database.PlayerInfo) error {

	// Data map to store [teamPrice][]teamPoints
	var m sync.Map

	// Simulate distribution strategy in batches (Prevent MySQL connection error 1040)
	for j := 0; j < (maxQueries / maxBatchSize); j++ {

		// Manage concurrency
		wg := &sync.WaitGroup{}
		wg.Add(maxBatchSize)

		for i := 0; i < maxBatchSize; i++ {
			go func() {
				defer wg.Done()

				team := <-simulatedTeams

				// Calculate the overall team points
				teamPoints, err := r.CalculateTeamPoints(team)
				if err != nil {
					fmt.Println(err)
				}

				// Calculate the overall team price
				teamPrice := CalculatePrice(team)

				// Add points value to correct map position, using sync map for concurrency
				currentPoints, ok := m.Load(teamPrice)
				if !ok {
					m.Store(teamPrice, []int{teamPoints})
				} else {
					currentPoints = append(currentPoints.([]int), teamPoints)
					m.Store(teamPrice, currentPoints)
				}
			}()
		}
		wg.Wait()
	}

	// For each possible map store, calculate the average
	consolidatedData := make([]string, 0)
	for i := 750; i <= 1000; i += 10 {

		// If the map position is empty, submit a zero value.
		// If the map position is not empty, calculate an average points value based on the contents of the map.
		if pointsTotal, ok := m.Load(i); !ok {
			consolidatedData = append(consolidatedData, fmt.Sprintf("%v, %v", i, 0))
		} else {
			pointsSum := 0
			for _, points := range pointsTotal.([]int) {
				pointsSum += points
			}
			averagePoints := float64(pointsSum) / float64(len(pointsTotal.([]int)))
			consolidatedData = append(consolidatedData, fmt.Sprintf("%v, %.2f", i, averagePoints))
		}
	}

	// Empty the results file of any old data, and write the new data to the file
	resultsFilePath := "internal/simulation_results/cost_variation.csv"
	if err := truncateFile(resultsFilePath); err != nil {
		return errors.Wrap(err)
	}
	if err := writeToFile(resultsFilePath, fmt.Sprintf("Team Price, Average Points\n")); err != nil {
		return errors.Wrap(err)
	}
	if err := writeToFile(resultsFilePath, strings.Join(consolidatedData, "\n")); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// RunDistributionStrategy simulates random teams, and records the points and cost distribution for each
func (r *Resolver) RunDistributionStrategy(simulatedTeams chan []database.PlayerInfo) error {

	// Data channels used to store simulation results
	resultsCh := make(chan []int, maxQueries)
	errCh := make(chan error, maxQueries)

	// Simulate distribution strategy in batches (Prevent MySQL connection error 1040)
	for j := 0; j < (maxQueries / maxBatchSize); j++ {

		// Manage concurrency
		wg := &sync.WaitGroup{}
		wg.Add(maxBatchSize)

		for i := 0; i < maxBatchSize; i++ {
			go func() {
				defer wg.Done()

				team := <-simulatedTeams

				// Calculate the overall team price.
				// If less than 950, ignore, since not using all available funds would skew the results
				if CalculatePrice(team) < 950 {
					return
				} else {
					fmt.Println(team)
				}

				// Calculate the overall team points
				teamPoints, err := r.CalculateTeamPoints(team)
				if err != nil {
					errCh <- err
				}

				// Calculate the cost distribution of the team, then add data into appropriate channels
				costDistribution, err := CalculateTeamDistribution(team)
				if err != nil {
					errCh <- err
				}

				resultsCh <- []int{costDistribution[0], teamPoints}
			}()
		}
		wg.Wait()
	}

	close(resultsCh)
	close(errCh)

	// Create an array to house each category of distribution, between 0 and 10
	distributionResults := make([][]int, 10)

	// For each result simulated, store result in the correct array space
	for result := range resultsCh {
		distributionResults[result[0]] = append(distributionResults[result[0]], result[1])
	}

	// TODO ...
	percentiles := make([][]int, 10)
	for key, category := range distributionResults {
		percentiles[key] = findPercentiles(category)
	}

	fmt.Println(percentiles[5])

	// Empty the results file of any old data, and write the new data to the file
	//resultsFilePath := "internal/simulation_results/cost_variation.csv"
	//if err := truncateFile(resultsFilePath); err != nil {
	//	return errors.Wrap(err)
	//}
	//if err := writeToFile(resultsFilePath, fmt.Sprintf("Team Price, Average Points\n")); err != nil {
	//	return errors.Wrap(err)
	//}
	//if err := writeToFile(resultsFilePath, strings.Join(consolidatedData, "\n")); err != nil {
	//	return errors.Wrap(err)
	//}

	return nil
}

// CalculateTeamDistribution takes the team and calculates what tier each player fits into
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
