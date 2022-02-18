package main

import (
	"errors"
	"math/rand"
	"time"
)

// FIXME Why are there 75 squares? Should be 25 squares, and we just keep track of what number
// each square is and which numbers have been called. If we don't have a square for a number,
// we don't care anything about it other than not calling it again.
// FIXME going to have a problem with the zero indexing

type Game struct {
	Squares         [25]Square   `json:"squares"`   //
	Called          map[int]bool `json:"called"`    // All the numbers that have been called
	PreCalled       map[int]bool `json:"preCalled"` // All the numbers that have been called
	Shapes          map[int]bool `json:"shapes"`    // Which shapes exist in the game
	BingoNumToIndex map[int]int
	rand            *rand.Rand `copier:"-"`
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
	g.BingoNumToIndex = make(map[int]int)
	g.Called[12] = true
	g.PreCalled[12] = true
	g.Shapes = make(map[int]bool)
	for i := 0; i < 25; i++ {
		g.BingoNumToIndex[i] = i
		s := newSquare(i)
		if i == 12 { // free square
			s.PreCalled = true
			s.Called = true
		}
		g.Squares[i] = s
	}
	// remove any mapping to the free square
	delete(g.BingoNumToIndex, 12)

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
func (g *Game) playSquare(bingoNum int) bool {
	// See if we have the square for the number that was called
	squareId, exists := g.BingoNumToIndex[bingoNum]
	if !exists { // we don't have that square
		return false
	}

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
