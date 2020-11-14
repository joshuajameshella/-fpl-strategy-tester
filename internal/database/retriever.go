package database

import (
	"math/rand"

	"github.com/doug-martin/goqu/v9"
	"github.com/icelolly/go-errors"
)

// PlayerInfo is the structure of data found in the GW1 table
type PlayerInfo struct {
	ID        int
	FirstName string
	LastName  string
	Position  string
	Price     int
	Team      int
}

// GetRandomPlayer searches the database for a random, cheap player
func (r *Resolver) GetRandomPlayer(position string) (PlayerInfo, error) {
	query, args, err := r.sqlBuilder.From(dataGW1).Where(
		goqu.And(
			goqu.C("position").Eq(position),
			goqu.C("price").Lte(50),
		),
	).ToSQL()

	if err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	rows, err := r.FPLDB.Query(query, args...)
	if err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	suitablePlayers := make([]PlayerInfo, 0)
	for rows.Next() {
		var player PlayerInfo
		if err := rows.Scan(
			&player.ID,
			&player.FirstName,
			&player.LastName,
			&player.Position,
			&player.Team,
			&player.Price,
		); err != nil {
			_ = rows.Close()
			return PlayerInfo{}, errors.Wrap(err)
		}
		suitablePlayers = append(suitablePlayers, player)
	}

	if err := rows.Close(); err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	if len(suitablePlayers) == 0 {
		return PlayerInfo{}, errors.New("Empty db response")
	}

	//rand.Seed(time.Now().UnixNano())

	// Return a random player from the list
	return suitablePlayers[rand.Intn(len(suitablePlayers))], nil
}

// UpgradePlayer takes the player passed in, and finds a more expensive alternative
func (r *Resolver) UpgradePlayer(player PlayerInfo) (PlayerInfo, error) {
	query, args, err := r.sqlBuilder.From(dataGW1).Where(
		goqu.And(
			goqu.C("position").Eq(player.Position),
			goqu.C("price").Gt(player.Price),
		),
	).ToSQL()

	if err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	rows, err := r.FPLDB.Query(query, args...)
	if err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	suitablePlayers := make([]PlayerInfo, 0)
	for rows.Next() {
		var player PlayerInfo
		if err := rows.Scan(
			&player.ID,
			&player.FirstName,
			&player.LastName,
			&player.Position,
			&player.Team,
			&player.Price,
		); err != nil {
			_ = rows.Close()
			return PlayerInfo{}, errors.Wrap(err)
		}
		suitablePlayers = append(suitablePlayers, player)
	}

	if err := rows.Close(); err != nil {
		return PlayerInfo{}, errors.Wrap(err)
	}

	// If no upgrade is available, return back the player in question
	if len(suitablePlayers) == 0 {
		return player, nil
	}

	// Return a random player from the list
	return suitablePlayers[rand.Intn(len(suitablePlayers))], nil
}
