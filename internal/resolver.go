package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math/rand"

	"github.com/icelolly/go-errors"
)

// Resolver is the entry-point for accessing the football data
type Resolver struct {
	Database *database.Resolver
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
	for calculatePrice(teamSelection) < maxValue {
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
