package gamemaster

import (
	"math/rand"
	"time"

	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
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

	rand *rand.Rand
}

func New(players []players.Player) (Gamemaster, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	state, err := rules.NewGame(len(players), r)
	return Gamemaster{
		Players:   players,
		Gamestate: state,
		Wins:      make([]int, state.NumPlayers),
		rand:      r,
	}, err
}

func (master *Gamemaster) TakeTurn() {
	action := master.Players[master.ActivePlayer].PlayCard(state.NewSimple(master.Gamestate))
	master.PlayCard(action, master.rand)
}

func (master *Gamemaster) PlayGame() {
	for !master.GameEnded {
		master.TakeTurn()
	}
	master.Wins[master.Winner] += 1
}

// PlaySeries plays an entire series with the provided players, returning the id of the player who won.
// A player wins after winning gamesToWin games.
func (master *Gamemaster) PlaySeries(gamesToWin int) (int, error) {
	for {
		pid, score, tie := master.HighScore()
		if score >= gamesToWin && !tie {
			return pid, nil
		}

		master.PlayGame()
		winner := (master.Winner - master.startPlayerOffset + master.NumPlayers) % master.NumPlayers
		master.startPlayerOffset = winner
		var err error
		master.Gamestate, err = rules.NewGame(master.NumPlayers, master.rand)
		if err != nil {
			return 0, err
		}
	}
}

// PlayStatistics plays totalGames with the players in a fixed order, returning the number of times player 0 won.
func (master *Gamemaster) PlayStatistics(totalGames int) (int, error) {
	for i := 0; i < totalGames; i++ {
		master.PlayGame()
		var err error
		master.Gamestate, err = rules.NewGame(master.NumPlayers, master.rand)
		if err != nil {
			return 0, err
		}
	}
	return master.Wins[0], nil
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
			tie = false
		} else if score == maxScore {
			tie = true
		}
	}
	return maxPid, maxScore, tie
}
