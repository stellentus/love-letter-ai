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

type WebPlayer struct {
	action chan rules.Action
}

func (wp WebPlayer) PlayCard(players.SimpleState) rules.Action {
	// Wait on the channel until a rules.Action is received, then return it.
	return <-wp.action
}

const NUMBER_OF_PLAYERS = 2

func main() {
	rand.Seed(time.Now().UnixNano())

	wp := WebPlayer{action: make(chan rules.Action, 1)}

	state, err := rules.NewGame(NUMBER_OF_PLAYERS)
	if err != nil {
		panic(err)
	}
	state.EventLog = rules.EventLog{PlayerNames: []string{"Human", "Computer"}}

	go func() {
		// The series can play in the background because it's mostly blocking for user input.
		// This doesn't shut down properly when the server shuts down.
		// If multiple users connect, bad things happen.
		for {
			var ply players.Player
			if state.ActivePlayer == 0 {
				ply = &wp
			} else {
				ply = &players.RandomPlayer{}
			}
			action := ply.PlayCard(players.NewSimpleState(state))
			state.PlayCard(action)

			if state.GameEnded {
				oldEL := state.EventLog
				state, err = rules.NewGame(NUMBER_OF_PLAYERS)
				if err != nil {
					panic(err)
				}
				state.EventLog = oldEL // Continue event log from before
			}
		}
	}()

	tmpl := template.Must(template.ParseFiles("../../res/templates/index.template.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// TODO parse the action, send it through the wp.chan, then get current game state and output it.
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

			// Play the action
			wp.action <- act
			for len(wp.action) == cap(wp.action) {
				// Wait until the channel has been read; then assume the action has been played
				time.Sleep(time.Millisecond)
			}

			// Now reload the content...
		}
		tmpl.Execute(w, stateForTemplate(state))
	})
	http.ListenAndServe(":8080", nil)
}

func stateForTemplate(state rules.Gamestate) LoveLetterState {
	fmt.Println(state)

	data := LoveLetterState{
		RevealedCards: state.Faceup.String(),
		Score: Score{
			You:      8,
			Computer: 2,
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
