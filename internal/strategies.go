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
	const maxQueries int = 500
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

	calculateAverages(resultsCh)

	return nil
}

func calculateAverages(resultsCh chan []int) {
	two := make([]int, 0)
	three := make([]int, 0)
	four := make([]int, 0)
	five := make([]int, 0)
	six := make([]int, 0)
	seven := make([]int, 0)

	twoSum := 0
	threeSum := 0
	fourSum := 0
	fiveSum := 0
	sixSum := 0
	sevenSum := 0

	for x := range resultsCh {
		if x[0] == 2 {
			two = append(two, x[1])
			twoSum += x[1]
		}
		if x[0] == 3 {
			three = append(three, x[1])
			threeSum += x[1]
		}
		if x[0] == 4 {
			four = append(four, x[1])
			fourSum += x[1]
		}
		if x[0] == 5 {
			five = append(five, x[1])
			fiveSum += x[1]
		}
		if x[0] == 6 {
			six = append(six, x[1])
			sixSum += x[1]
		}
		if x[0] == 7 {
			seven = append(seven, x[1])
			sevenSum += x[1]
		}
	}

	fmt.Printf("Two Average: %v\n", float64(twoSum)/float64(len(two)))
	fmt.Printf("Three Average: %v\n", float64(threeSum)/float64(len(three)))
	fmt.Printf("Four Average: %v\n", float64(fourSum)/float64(len(four)))
	fmt.Printf("Five Average: %v\n", float64(fiveSum)/float64(len(five)))
	fmt.Printf("Six Average: %v\n", float64(sixSum)/float64(len(six)))
	fmt.Printf("Seven Average: %v\n", float64(sevenSum)/float64(len(seven)))
}
