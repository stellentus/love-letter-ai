package rules

import "math/rand"

type Card int

// Stack is an ordered stack of cards.
type Stack []Card

// For now, the Card value is its face value.
const (
	None = Card(iota) // Used to indicate errors or other things
	Guard
	Priest
	Baron
	Handmaid
	Prince
	King
	Countess
	Princess
)

// Deck contains Cards with an integer count
type Deck map[Card]int

func DefaultDeck() Deck {
	return Deck{
		Guard:    5,
		Priest:   2,
		Baron:    2,
		Handmaid: 2,
		Prince:   2,
		King:     1,
		Countess: 1,
		Princess: 1,
	}
}

func CardNames() []Card {
	return []Card{
		Guard,
		Priest,
		Baron,
		Handmaid,
		Prince,
		King,
		Countess,
		Princess,
	}
}

func (deck Deck) Size() int {
	sum := 0
	for _, val := range deck {
		sum += val
	}
	return sum
}

func (deck Deck) CountFor(card Card) int {
	return map[Card]int(deck)[card] // Possibly the most unreadable line of go I've ever written :)
}

func (deck *Deck) Draw() Card {
	draw := int(rand.Int31n(int32(deck.Size())))
	for _, name := range CardNames() {
		thisCount := deck.CountFor(name)
		if draw < thisCount {
			map[Card]int(*deck)[name] -= 1 // Or maybe this is even more unreadable
			return name
		} else {
			draw -= thisCount
		}
	}
	return None
}

func (deck *Deck) AddStack(stack Stack) {
	for _, card := range stack {
		map[Card]int(*deck)[card] += 1
	}
}

func (stack Stack) AsDeck() Deck {
	deck := map[Card]int{}
	for _, name := range CardNames() {
		deck[name] = 0
	}
	for _, card := range stack {
		deck[card] += 1
	}
	return deck
}

func (stack Stack) Score() int {
	sum := 0
	for _, card := range stack {
		sum += int(card)
	}
	return sum
}
