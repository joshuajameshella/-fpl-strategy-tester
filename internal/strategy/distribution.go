package strategy

import (
	"fpl-strategy-tester/internal/database"
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
