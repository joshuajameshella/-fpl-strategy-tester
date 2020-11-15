package database

import (
	"math/rand"

	"github.com/doug-martin/goqu/v9"
	"github.com/icelolly/go-errors"
)

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

// GetPlayerData takes the player ID and returns the data for each match played
func (r *Resolver) GetPlayerData(playerID int) ([]PlayerGWInfo, error) {
	query, args, err := r.sqlBuilder.From(playerData).Where(
		goqu.C("element").Eq(playerID),
	).ToSQL()

	if err != nil {
		return nil, errors.Wrap(err)
	}

	rows, err := r.FPLDB.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	playerData := make([]PlayerGWInfo, 0)
	for rows.Next() {
		var gw PlayerGWInfo
		if err := rows.Scan(
			&gw.Name,
			&gw.Element,
			&gw.OpponentTeam,
			&gw.TotalPoints,
			&gw.Value,
			&gw.WasHome,
			&gw.GW,
		); err != nil {
			_ = rows.Close()
			return nil, errors.Wrap(err)
		}
		playerData = append(playerData, gw)
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Wrap(err)
	}

	if len(playerData) == 0 {
		return nil, errors.New("Empty db response")
	}

	// Return a random player from the list
	return playerData, nil
}
