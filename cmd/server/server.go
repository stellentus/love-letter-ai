package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"love-letter-ai/montecarlo"
	"love-letter-ai/players"
	"love-letter-ai/rules"
	"love-letter-ai/state"
	"love-letter-ai/td"

	"github.com/kelseyhightower/envconfig"
)

const NUMBER_OF_PLAYERS = 2

var (
	sarsaFile = flag.String("sarsa", "", "Path to a sarsa file")
	qFile     = flag.String("q", "", "Path to a Q learning file")

	config = struct {
		Resources string `default:"../../res"`
		Address   string `default:":8080"`
	}{}
)

func exitIfError(err error, reason string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "exiting: %s\n%s\n", reason, err)
		os.Exit(1)
	}
}

func resourcePath(path string) string {
	return filepath.Join(config.Resources, path)
}

type GameData struct {
	Game, Score, Opponents interface{}
}

func main() {
	flag.Parse()
	if *sarsaFile != "" && *qFile != "" {
		exitIfError(errors.New("Can only specify one of -sarsa or -q"), "invalid arguments")
	}

	exitIfError(envconfig.Process("LLAI", &config), "failed to parse environment")

	bots := map[string]players.Player{
		"random": &players.RandomPlayer{},
	}

	if *sarsaFile != "" {
		sarsa := td.NewTD(0, 0)
		exitIfError(sarsa.LoadFromFile(*sarsaFile), "loading sarsa file")
		bots["sarsa"] = sarsa.SarsaLearner()
	}

	if *qFile != "" {
		q := montecarlo.NewQPlayer(0)
		exitIfError(q.LoadFromFile(*qFile), "loading Q file")
		bots["q"] = q
	}

	rand.Seed(time.Now().UnixNano())

	score := []int{0, 0} // Number of wins for each player

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(resourcePath("static")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		game := gameFromToken(cookieID(r))

		botName := "random"

		switch r.Method {
		case "POST":
			rand := rand.New(rand.NewSource(time.Now().UnixNano()))
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
			game.PlayCard(act, rand)

			// Did the player's move end the game?
			if game.GameEnded {
				score[game.Winner]++
				game.Reset(rand)
				break
			}

			botName = strings.ToLower(r.FormValue("opponent"))
			comPlay, ok := bots[botName]
			if !ok {
				log.Printf("Invalid botname (%s), defaulting to random", botName)
				botName = "random"
				comPlay = bots[botName]
			}

			// The player didn't end the game, so the computer gets a turn...
			action := comPlay.PlayCard(state.NewSimple(game))
			game.PlayCard(action, rand)

			if game.GameEnded {
				score[game.Winner]++
				game.Reset(rand)
			}

			// Now reload the content...
		}

		gd := GameData{
			Game: stateForTemplate(game),
			Score: struct {
				You      int
				Computer int
			}{
				You:      score[0],
				Computer: score[1],
			},
			Opponents: Opponents(botName, bots),
		}

		err := template.Must(template.ParseFiles(resourcePath("templates/index.template.html"))).Execute(w, gd)
		if err != nil {
			fmt.Println("Error:", err)
		}
	})

	log.Println("Running server at", config.Address)

	http.ListenAndServe(config.Address, nil)
}

func Opponents(current string, bots map[string]players.Player) interface{} {
	list := []string{}
	for o := range bots {
		list = append(list, strings.ToTitle(o))
	}
	return struct {
		Current string
		Bots    []string
	}{
		strings.ToTitle(current),
		list,
	}
}

func cookieID(r *http.Request) string {
	cookie, err := r.Cookie("GameStateID")
	if err != nil {
		return ""
	}
	data, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return ""
	}
	return data
}

func gameFromToken(tok string) rules.Gamestate {
	var game rules.Gamestate
	err := errors.New("")

	if tok != "" {
		// Try loading if there's a token
		err = game.FromToken(tok)
	}

	if err != nil {
		game, err = rules.NewGame(NUMBER_OF_PLAYERS, rand.New(rand.NewSource(time.Now().UnixNano())))
		if err != nil {
			panic(err)
		}
	}

	game.EventLog = rules.EventLog{PlayerNames: []string{"Human", "Computer"}}
	return game
}

func stateForTemplate(game rules.Gamestate) interface{} {
	type PlayedCards struct {
		You      []string
		Computer []string
	}

	type LoveLetterState struct {
		RevealedCards []string
		PlayedCards
		LastPlay    string
		Card1       string
		Card2       string
		EventLog    template.HTML
		GameStateID string
	}

	data := LoveLetterState{
		RevealedCards: game.Faceup.Strings(),
		PlayedCards: PlayedCards{
			You:      game.Discards[0].Strings(),
			Computer: game.Discards[1].Strings(),
		},
		LastPlay:    game.LastPlay[1].String(),
		Card1:       game.CardInHand[0].String(),
		Card2:       game.ActivePlayerCard.String(), // TODO this assumes that the current player is the active player
		EventLog:    template.HTML(strings.Join(game.EventLog.Events, "<br>")),
		GameStateID: game.Token(),
	}
	return data
}
