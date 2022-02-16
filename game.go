package main

import (
	"errors"
	"math/rand"
	"time"
)

type Game struct {
	Squares   [76]Square   `json:"squares"`   //
	Called    map[int]bool `json:"called"`    // All the numbers that have been called
	PreCalled map[int]bool `json:"preCalled"` // All the numbers that have been called
	Shapes    map[int]bool `json:"shapes"`    // Which shapes exist in the game
	rand      *rand.Rand   `copier:"-"`
}

type Square struct {
	Number    int          `json:"number"`
	Called    bool         `json:"called"`
	PreCalled bool         `json:"preCalled"`
	Needed    map[int]bool `json:"needed"` // to handle multiple shapes
}

func newGame() *Game {
	g := Game{}
	g.Called = make(map[int]bool)
	g.PreCalled = make(map[int]bool)
	g.Called[13] = true
	g.PreCalled[13] = true
	g.Shapes = make(map[int]bool)
	for i := 1; i <= 75; i++ {
		s := newSquare(i)
		if i == 13 { // free square
			s.PreCalled = true
			s.Called = true
		}
		g.Squares[i] = s
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(s1)

	return &g
}

func (g *Game) setNeeded(shapeId int, squares []int) {
	g.Shapes[shapeId] = true
	for _, square := range squares {
		g.Squares[square].Needed[shapeId] = true
	}
}

// Checks if playing this square has created a win
func (g *Game) playSquare(squareId int) bool {
	g.Squares[squareId].Called = true
	for shapeId, _ := range g.Shapes {
		// doing it this way is slightly more efficient but only tells you if that move created a win,
		// not if the board was already in win state
		//for shapeId, _ := range g.Squares[squareId].Needed {
		if g.winner(shapeId) {
			return true
		}
	}
	return false
}

func (g *Game) winner(shape int) bool {
	for _, square := range g.Squares {
		if square.Needed[shape] && !square.Called {
			return false
		}
	}
	return true
}

func (g *Game) callRandomSquare() (int, error) {
	n := 0
	if len(g.Called) >= 75 {
		return 0, errors.New("No more numbers to pick")
	}
	for ; n == 0 || (n > 0 && g.Called[n] == true); n = g.rand.Intn(75) + 1 {
	}
	g.Called[n] = true
	return n, nil
}

// Undo all calls
func (g *Game) reset() {
	g.Called = make(map[int]bool)
	for key, value := range g.PreCalled {
		g.Called[key] = value
	}
	for i, square := range g.Squares {
		square.Called = square.PreCalled
		g.Squares[i] = square
	}
}

func newSquare(n int) Square {
	return Square{Number: n, Called: false, Needed: make(map[int]bool)}
}
