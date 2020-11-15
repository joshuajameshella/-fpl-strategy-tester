package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"fpl-strategy-tester/internal/strategy"
	"math/rand"
	"sync"

	"github.com/icelolly/go-errors"
)

// PickRandomTeam creates a random team from the player selections available in GW1
func (r *Resolver) PickRandomTeam() ([]database.PlayerInfo, error) {

	// Create an empty team
	teamSelection := make([]database.PlayerInfo, 0)

	// Select and add two random goalkeepers to the team
	for i := 0; i < 2; i++ {
		player, err := r.Database.GetRandomPlayer("G")
		if err != nil {
			return nil, errors.Wrap(err)
		}
		teamSelection = append(teamSelection, player)
	}

	// Select and add five random defenders to the team
	for i := 0; i < 5; i++ {
		player, err := r.Database.GetRandomPlayer("D")
		if err != nil {
			return nil, errors.Wrap(err)
		}
		teamSelection = append(teamSelection, player)
	}

	// Select and add five random midfielders to the team
	for i := 0; i < 5; i++ {
		player, err := r.Database.GetRandomPlayer("M")
		if err != nil {
			return nil, errors.Wrap(err)
		}
		teamSelection = append(teamSelection, player)
	}

	// Select and add three random forwards to the team
	for i := 0; i < 3; i++ {
		player, err := r.Database.GetRandomPlayer("F")
		if err != nil {
			return nil, errors.Wrap(err)
		}
		teamSelection = append(teamSelection, player)
	}

	// While the team's value remains under £95M, continue to upgrade random players in the team
	for calculatePrice(teamSelection) < 950 {
		randomPlayer := rand.Intn(len(teamSelection))

		if playerUpgrade, err := r.Database.UpgradePlayer(teamSelection[randomPlayer]); err != nil {
			fmt.Printf("Error occured while attempting to upgrade player: %v\n", err)
		} else {
			teamSelection[randomPlayer] = playerUpgrade
		}
	}

	return teamSelection, nil
}

// CalculateTeamPoints takes the team of players and returns a total of their end-of-season points
func (r *Resolver) CalculateTeamPoints(team []database.PlayerInfo) (int, error) {

	teamPoints := 0
	for _, player := range team {

		// Get each week of data for the specified player
		gwData, err := r.Database.GetPlayerData(player.ID)
		if err != nil {
			return 0, errors.Wrap(err)
		}

		// Add the game-week points to the tally
		for _, gw := range gwData {
			teamPoints += gw.TotalPoints
		}
	}

	return teamPoints, nil
}

// calculatePrice takes the team info and return's the combined worth of all players
func calculatePrice(team []database.PlayerInfo) int {
	teamPrice := 0
	for _, player := range team {
		teamPrice += player.Price
	}
	return teamPrice
}

// runDistributionStrategy simulates random teams, and records the points and cost distribution for each
func (r *Resolver) RunDistributionStrategy() error {

	// How many simulations to run
	const maxQueries int = 1000
	const maxBatchSize int = 100

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

				// Simulate a random FPL team
				team, err := r.PickRandomTeam()
				if err != nil {
					errCh <- err
				}

				// Calculate the overall team points
				points, err := r.CalculateTeamPoints(team)
				if err != nil {
					errCh <- err
				}

				// Calculate the cost distribution of the team
				costDistribution := strategy.CalculateTeamDistribution(team)

				// Add data into channels
				resultsCh <- []int{costDistribution[0], points}
				errCh <- err
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

	// Handle the simulation results
	strategy.ProcessDistributionResults(resultsCh)

	return nil
}
