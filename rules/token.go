package rules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const tokenFormat = "%d.%d.%s.%s.%s.%s.%d.%s.%s.%d"

func (game *Gamestate) Token() string {
	return fmt.Sprintf(tokenFormat,
		game.NumPlayers,
		game.Deck.AsInt(),
		game.Faceup.Token(),
		game.Discards.Token(),
		game.LastPlay.Token(),
		game.KnownCards.Token(),
		game.ActivePlayer,
		stringForElim(game.EliminatedPlayers),
		game.CardInHand.Token(),
		int(game.ActivePlayerCard))
}

func (game *Gamestate) FromToken(tok string) error {
	strs := strings.Split(tok, ".")
	if len(strs) != 10 {
		return errors.New("Token '" + tok + "' does not have the expected 10 values")
	}

	game.NumPlayers = unsafeAtoi(strs[0])
	game.Deck.FromInt(unsafeAtoi(strs[1]))
	game.Faceup.FromToken(strs[2])
	game.Discards.FromToken(strs[3])
	game.LastPlay.FromToken(strs[4])
	game.KnownCards.FromToken(strs[5])
	game.ActivePlayer = unsafeAtoi(strs[6])
	game.EliminatedPlayers = elimFromString(strs[7])
	game.CardInHand.FromToken(strs[8])
	game.ActivePlayerCard = Card(unsafeAtoi(strs[9]))

	game.EventLog = newEventLog(game.NumPlayers)
	game.EventLog.log(fmt.Sprintf("Game was created from token %s", tok))

	var err error
	if !game.isValid() {
		return errors.New("Token '" + tok + "' is not a valid game")
	}

	return err
}

func stringForElim(elim []bool) string {
	str := ""
	for _, elim := range elim {
		if elim {
			str += "1"
		} else {
			str += "0"
		}
	}
	return str
}
func elimFromString(str string) []bool {
	ep := []bool{}
	for _, c := range str {
		ep = append(ep, c == '1')
	}
	return ep
}

func unsafeAtoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

func (game *Gamestate) isValid() bool {
	if len(game.Discards) != game.NumPlayers {
		return false
	}
	if len(game.LastPlay) != game.NumPlayers {
		return false
	}
	if len(game.KnownCards) != game.NumPlayers {
		return false
	}
	if len(game.EliminatedPlayers) != game.NumPlayers {
		return false
	}
	if len(game.CardInHand) != game.NumPlayers {
		return false
	}
	if len(game.EventLog.PlayerNames) != game.NumPlayers {
		return false
	}
	return true
}
