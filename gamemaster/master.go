package gamemaster

import (
	"love-letter-ai/players"
	"love-letter-ai/rules"
)

type Gamemaster struct {
	// Players is a list of the players in the current game
	Players []players.Player

	// Gamestate tracks the state of the game
	rules.Gamestate

	// Wins tracks each player's wins so far
	Wins []int

	// startPlayerOffset is the id of the player who started the current game
	startPlayerOffset int
}

func New(players []players.Player) (Gamemaster, error) {
	state, err := rules.NewGame(len(players))
	return Gamemaster{
		Players:   players,
		Gamestate: state,
		Wins:      make([]int, state.NumPlayers),
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
	master.Wins[master.Winner] += 1
	return nil
}

// PlaySeries playes an entire series with the provided players, returning the id of the player who won.
// A player wins after winning gamesToWin games.
func (master *Gamemaster) PlaySeries(gamesToWin int) (int, error) {
	for {
		pid, score, tie := master.HighScore()
		if score >= gamesToWin && !tie {
			return pid, nil
		}

		err := master.PlayGame()
		if err != nil {
			return 0, err
		}
		winner := (master.Winner - master.startPlayerOffset) % master.NumPlayers
		master.Wins[winner] += 1
		master.startPlayerOffset = winner
		master.Gamestate, err = rules.NewGame(master.NumPlayers)
	}
}

// HighScore returns the player who scored highest and that player's score.
// It also returns a bool to indicate if there's a tie. There isn't currently any way to tie.
func (master *Gamemaster) HighScore() (int, int, bool) {
	maxPid := -1
	maxScore := 0
	tie := false
	for pid, score := range master.Wins {
		if score > maxScore {
			maxPid = pid
			maxScore = score
		} else if score == maxScore {
			tie = true
		}
	}
	return maxPid, maxScore, tie
}
