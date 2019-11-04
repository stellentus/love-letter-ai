package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"love-letter-ai/td"
)

type Score struct {
	You      int
	Computer int
}

type PlayedCards struct {
	You      []string
	Computer []string
}

type LoveLetterState struct {
	RevealedCards []string
	Score
	PlayedCards
	LastPlay    string
	Card1       string
	Card2       string
	EventLog    template.HTML
	GameStateID int
}

const NUMBER_OF_PLAYERS = 2

var (
	sarsaFile = flag.String("sarsa", "", "Path to a sarsa file")
	qFile     = flag.String("q", "", "Path to a Q learning file")
)

func exitIfError(err error, reason string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "exiting: %s\n%s\n", reason, err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	if *sarsaFile != "" && *qFile != "" {
		exitIfError(errors.New("Can only specify one of -sarsa or -q"), "invalid arguments")
	}

	var comPlay players.Player
	switch {
	case *sarsaFile != "":
		sarsa := td.NewSarsa(0, 0, 0)
		exitIfError(sarsa.LoadFromFile(*sarsaFile), "loading sarsa file")
		comPlay = sarsa.NewPlayer()
	case *qFile != "":
		q := montecarlo.NewQPlayer(0)
		exitIfError(q.LoadFromFile(*qFile), "loading Q file")
		comPlay = q
	default:
		comPlay = &players.RandomPlayer{}
	}

	rand.Seed(time.Now().UnixNano())

	score := []int{0, 0} // Number of wins for each player

	game, err := rules.NewGame(NUMBER_OF_PLAYERS)
	if err != nil {
		panic(err)
	}
	game.EventLog = rules.EventLog{PlayerNames: []string{"Human", "Computer"}}

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("../../res/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if game.ActivePlayer != 0 {
				// This is an error
				panic("Active player is not the human!")
			}

			if err := r.ParseForm(); err != nil {
				fmt.Printf("ParseForm() err: %v", err)
				return
			}

			// Parse the action
			act := rules.Action{}
			if r.FormValue("cards") == "card2" {
				act.PlayRecent = true
			}
			if r.FormValue("targets") == "computer" {
				act.TargetPlayerOffset = 1
			}
			act.SelectedCard = rules.CardFromString(r.FormValue("guess"))
			fmt.Println("Parsing action:", act)
			game.PlayCard(act)

			// Did the player's move end the game?
			if game.GameEnded {
				score[game.Winner]++
				game.Reset()
				break
			}

			// The player didn't end the game, so the computer gets a turn...
			action := comPlay.PlayCard(state.NewSimple(game))
			game.PlayCard(action)

			if game.GameEnded {
				score[game.Winner]++
				game.Reset()
			}

			// Now reload the content...
		}
		err := template.Must(template.ParseFiles("../../res/templates/index.template.html")).Execute(w, stateForTemplate(game, score))
		if err != nil {
			fmt.Println("Error:", err)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func stateForTemplate(game rules.Gamestate, score []int) LoveLetterState {
	fmt.Println(game)

	data := LoveLetterState{
		RevealedCards: game.Faceup.Strings(),
		Score: Score{
			You:      score[0],
			Computer: score[1],
		},
		PlayedCards: PlayedCards{
			You:      game.Discards[0].Strings(),
			Computer: game.Discards[1].Strings(),
		},
		LastPlay:    game.LastPlay[1].String(),
		Card1:       game.CardInHand[0].String(),
		Card2:       game.ActivePlayerCard.String(), // TODO this assumes that the current player is the active player
		EventLog:    template.HTML(strings.Join(game.EventLog.Events, "<br>")),
		GameStateID: state.NewSimple(game).AsIndex(),
	}
	return data
}
