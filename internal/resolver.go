package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/icelolly/go-errors"
	"github.com/patrickmn/go-cache"
)

// How many simulations to run
const maxQueries int = 10000
const maxBatchSize int = 50

// Resolver is the entry-point for accessing the football data
type Resolver struct {
	Database *database.Resolver
	Cache    *cache.Cache
}

// NewResolver creates and returns an empty Resolver
func NewResolver() *Resolver {
	return &Resolver{}
}

// ResolveDatabase returns or initiates a new database connection
func (r *Resolver) ResolveDatabase() *database.Resolver {
	if r.Database == nil {
		repo := database.NewResolver()
		repo.ResolveFPLDB()
		repo.ResolveMySQLQueryBuilder()
		r.Database = repo
	}
	return r.Database
}

// ResolveCache creates a new, or re-uses an existing cache instance
func (r *Resolver) ResolveCache() *cache.Cache {
	if r.Cache == nil {
		r.Cache = cache.New(5*time.Minute, 10*time.Minute)
	}
	return r.Cache
}

// GenerateTeams simulates 10,000 possible teams, and returns them on a channel for the results
// to be analysed by the different strategies.
func (r *Resolver) GenerateTeams() (chan []database.PlayerInfo, chan error) {

	// Data channels used to store simulation simulation_results
	resultsCh := make(chan []database.PlayerInfo, maxQueries)
	errCh := make(chan error, maxQueries)

	// Simulate distribution strategy in batches (Prevent MySQL connection error 1040)
	for j := 0; j < (maxQueries / maxBatchSize); j++ {

		// Manage concurrency
		wg := &sync.WaitGroup{}
		wg.Add(maxBatchSize)

		for i := 0; i < maxBatchSize; i++ {
			go func() {
				defer wg.Done()

				// Create a random team value to simulate (between £75M & £100M)
				randomTeamValue := rand.Intn(100-75) + 75

				// Simulate a random FPL team, up to the maximum value
				if team, err := r.PickRandomTeam(randomTeamValue * 10); err != nil {
					errCh <- err
				} else {
					resultsCh <- team
				}
			}()
		}
		wg.Wait()
	}

	close(errCh)
	close(resultsCh)

	return resultsCh, errCh
}

// PickRandomTeam creates a random team from the player selections available in GW1
// It takes the maximum value a team can be, and returns a team equal to that value
func (r *Resolver) PickRandomTeam(maxValue int) ([]database.PlayerInfo, error) {

	// The minimum value a team can be is £75M
	if maxValue < 750 {
		return nil, errors.New("Unable to create a team - team value too low")
	}

	// Create an empty team
	teamSelection := make([]database.PlayerInfo, 0)

	// Select and add two random goalkeepers to the team
	for i := 0; i < 2; i++ {
		selectedPlayer, err := r.Database.GetRandomPlayer("G")
		if err != nil {
			return nil, errors.Wrap(err)
		}
		teamSelection = append(teamSelection, selectedPlayer)
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

	// While the team's value remains under the maximum value, continue to upgrade random players in the team
	for CalculatePrice(teamSelection) < maxValue {
		randomPlayer := rand.Intn(len(teamSelection))

		if playerUpgrade, err := r.Database.UpgradePlayer(teamSelection[randomPlayer]); err != nil {
			fmt.Printf("Error occured while attempting to upgrade player: %v\n", err)
		} else {
			teamSelection[randomPlayer] = playerUpgrade
		}
	}

	// If the team's value exceeds the £100M budget, continue to downgrade random players in the team
	for CalculatePrice(teamSelection) > 1000 {
		randomPlayer := rand.Intn(len(teamSelection))

		if playerDowngrade, err := r.Database.DowngradePlayer(teamSelection[randomPlayer]); err != nil {
			fmt.Printf("Error occured while attempting to upgrade player: %v\n", err)
		} else {
			teamSelection[randomPlayer] = playerDowngrade
		}
	}

	// Cycle through each player in the team to check for any rules broken
	teamCriteria := false
	for !teamCriteria {
		duplicatePlayers := false

		// Check that there are no duplicate players in the team
		for key, player := range teamSelection {
			for i := key + 1; i < len(teamSelection); i++ {
				if player.ID == teamSelection[i].ID {
					duplicatePlayers = true

					// Replace the duplicate player with an equivalent alternative
					replacementPlayer, err := r.Database.ReplacePlayer(player, teamSelection)
					if err != nil {
						fmt.Printf("Error while replacing player: %v\n", err)
						continue
					}
					teamSelection[i] = replacementPlayer
				}
			}
		}
		if !duplicatePlayers {
			teamCriteria = true
		}
	}

	return teamSelection, nil
}

// CalculateTeamPoints takes the team of players and returns a total of their end-of-season points
func (r *Resolver) CalculateTeamPoints(team []database.PlayerInfo) (int, error) {

	teamPoints := 0
	for _, player := range team {

		// Check if the player points have already been calculated
		playerPoints, found := r.Cache.Get(strconv.Itoa(player.ID))

		// If the player data has not yet been calculated, perform the calculation and store in cache
		if !found {
			// Get the 38 weeks of player data
			gwData, err := r.Database.GetPlayerData(player.ID)
			if err != nil {
				return 0, errors.Wrap(err)
			}

			pointsTotal := 0
			for _, gw := range gwData {
				pointsTotal += gw.TotalPoints
			}
			playerPoints = pointsTotal

			// Save the value in cache
			r.Cache.Set(strconv.Itoa(player.ID), pointsTotal, cache.DefaultExpiration)
		}
		teamPoints += playerPoints.(int)
	}

	return teamPoints, nil
}

// CalculatePrice takes the team info and return's the combined worth of all players
func CalculatePrice(team []database.PlayerInfo) int {
	teamPrice := 0
	for _, player := range team {
		teamPrice += player.Price
	}
	return teamPrice
}

// truncateFile empties the desired file, ready for new data
func truncateFile(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrap(err)
	}
	if err = f.Close(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func writeToFile(filePath, textBody string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err)
	}
	defer f.Close()
	if _, err := f.WriteString(textBody); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
