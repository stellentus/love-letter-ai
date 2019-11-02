package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"love-letter-ai/players"
	"love-letter-ai/rules"
)

type Score struct {
	You      int
	Computer int
}

type PlayedCards struct {
	You      string
	Computer string
}

type LoveLetterState struct {
	RevealedCards string
	Score
	PlayedCards
	LastPlay string
	Card1    string
	Card2    string
	EventLog template.HTML
}

const NUMBER_OF_PLAYERS = 2

func main() {
	rand.Seed(time.Now().UnixNano())

	comPlay := &players.RandomPlayer{}
	score := []int{0, 0} // Number of wins for each player

	state, err := rules.NewGame(NUMBER_OF_PLAYERS)
	if err != nil {
		panic(err)
	}
	state.EventLog = rules.EventLog{PlayerNames: []string{"Human", "Computer"}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if state.ActivePlayer != 0 {
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
			state.PlayCard(act)

			// Did the player's move end the game?
			if state.GameEnded {
				score[state.Winner]++
				state.Reset()
				break
			}

			// The player didn't end the game, so the computer gets a turn...
			action := comPlay.PlayCard(players.NewSimpleState(state))
			state.PlayCard(action)

			if state.GameEnded {
				score[state.Winner]++
				state.Reset()
			}

			// Now reload the content...
		}
		err := template.Must(template.ParseFiles("../../res/templates/index.template.html")).Execute(w, stateForTemplate(state, score))
		if err != nil {
			fmt.Println("Error:", err)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func stateForTemplate(state rules.Gamestate, score []int) LoveLetterState {
	fmt.Println(state)

	data := LoveLetterState{
		RevealedCards: state.Faceup.String(),
		Score: Score{
			You:      score[0],
			Computer: score[1],
		},
		PlayedCards: PlayedCards{
			You:      state.Discards[0].String(),
			Computer: state.Discards[1].String(),
		},
		LastPlay: state.LastPlay[1].String(),
		Card1:    state.CardInHand[0].String(),
		Card2:    state.ActivePlayerCard.String(), // TODO this assumes that the current player is the active player
		EventLog: template.HTML(strings.Join(state.EventLog.Events, "<br>")),
	}
	return data
}
