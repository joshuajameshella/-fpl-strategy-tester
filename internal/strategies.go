package internal

import (
	"fmt"
	"fpl-strategy-tester/internal/database"
	"math/rand"

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

// calculatePrice takes the team info and return's the combined worth of all players
func calculatePrice(team []database.PlayerInfo) int {
	teamPrice := 0
	for _, player := range team {
		teamPrice += player.Price
	}
	return teamPrice
}
