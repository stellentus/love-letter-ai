package main

import (
	"fmt"
	"html/template"
	"net/http"
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
}

func main() {
	tmpl := template.Must(template.ParseFiles("../../res/templates/index.template.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := LoveLetterState{
			RevealedCards: "Guard, Baron, Handmaid",
			Score: Score{
				You:      8,
				Computer: 2,
			},
			PlayedCards: PlayedCards{
				You:      "Priest, Prince, Baron, Handmaid",
				Computer: "Guard, Countess, Guard",
			},
			LastPlay: "Guard, guessing you had a Princess",
			Card1:    "Priest",
			Card2:    "Guard",
		}

		fmt.Println("Received", r.Method)

		switch r.Method {
		case "GET":
			tmpl.Execute(w, data)
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Printf("ParseForm() err: %v", err)
				return
			}
			fmt.Printf("cards = %s\n", r.FormValue("cards"))
			fmt.Printf("targets = %s\n", r.FormValue("targets"))
			fmt.Printf("guess = %s\n", r.FormValue("guess"))

			tmpl.Execute(w, data)
		}
	})
	http.ListenAndServe(":8080", nil)
}
