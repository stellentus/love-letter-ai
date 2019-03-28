package gamemaster

import (
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

type Gamemaster struct {
	Players []players.Player
	rules.Gamestate
}

func New(players []players.Player) (Gamemaster, error) {
	state, err := rules.NewGame(len(players))
	return Gamemaster{
		Players:   players,
		Gamestate: state,
	}, err
}

func (master *Gamemaster) TakeTurn() error {
	action := master.Players[master.ActivePlayer].PlayCard(master.Gamestate.AsSimpleState(), master.ActivePlayer)
	return master.PlayCard(action)
}

func (master *Gamemaster) PlayGame() error {
	for !master.GameEnded {
		if err := master.TakeTurn(); err != nil {
			return err
		}
	}
	return nil
}
