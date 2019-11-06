package rules

import (
	"math/rand"
	"strings"
)

type Card int

// Stack is an ordered stack of cards.
type Stack []Card

func (st Stack) String() string {
	return strings.Join(st.Strings(), ", ")
}

func (st Stack) Strings() []string {
	strs := make([]string, len(st))
	for i, c := range st {
		strs[i] = c.String()
	}
	return strs
}

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
	numberOfCards
)

var nameOfCard = map[Card]string{
	Guard:    "Guard",
	Priest:   "Priest",
	Baron:    "Baron",
	Handmaid: "Handmaid",
	Prince:   "Prince",
	King:     "King",
	Countess: "Countess",
	Princess: "Princess",
}

func (c Card) String() string {
	str, _ := nameOfCard[c]
	return str
}

func CardFromString(s string) Card {
	for c, str := range nameOfCard {
		if strings.EqualFold(str, s) {
			return c
		}
	}
	return None
}

// Deck contains counts of Cards
type Deck [numberOfCards]int

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

func (deck Deck) Copy() Deck {
	deck2 := Deck{}
	for i, val := range deck {
		deck2[i] = val
	}
	return deck2
}

func (deck *Deck) Draw() Card {
	size := deck.Size()
	if size == 0 {
		return None
	}
	draw := int(rand.Int31n(int32(size)))
	for name, count := range deck {
		if draw < count {
			deck[name] -= 1
			return Card(name)
		} else {
			draw -= count
		}
	}
	return None
}

func (deck *Deck) AddStack(stack Stack) {
	for _, card := range stack {
		deck[card] += 1
	}
}

func (stack Stack) AsDeck() Deck {
	deck := Deck{}
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

func (stack Stack) Copy() Stack {
	stack2 := make([]Card, 0, len(stack))
	for _, val := range stack {
		stack2 = append(stack2, val)
	}
	return stack2
}

const DeckSpaceMagnitude = 6 * 3 * 3 * 3 * 3 * 2 * 2 * 2

func (sc Deck) AsInt() int {
	// uses 12 bits
	return sc[Guard] +
		6*(sc[Priest]+
			3*(sc[Baron]+
				3*(sc[Handmaid]+
					3*(sc[Prince]+
						3*(sc[Princess]+
							2*(sc[Countess]+
								2*sc[King]))))))
}

func (sc *Deck) FromInt(i int) {
	i, sc[Guard] = divRem(i, 6)
	i, sc[Priest] = divRem(i, 3)
	i, sc[Baron] = divRem(i, 3)
	i, sc[Handmaid] = divRem(i, 3)
	i, sc[Prince] = divRem(i, 3)
	i, sc[Princess] = divRem(i, 2)
	i, sc[Countess] = divRem(i, 2)
	sc[King] = int(i)
}
func divRem(num int, den int) (int, int) { return num / den, num % den }

type Stacks []Stack
