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

/*
so a Game contains a map of bingo numbers to Squares

when we create a Game, pass in Squares
each Square knows its number, and for which Shapes it is needed
the Game knows the list of Squares, and which numbers are called.
does it need to know the shapes? i think so, just so we know which shapes to query

*/

// These methods ignore the free square- when you're creating the board just mark it called on
// creation, or mark it as not required.

type Game struct {
	Squares     map[int]Square `json:"squares"`   // The squares on the player's board
	Called      map[int]bool   `json:"called"`    // All the numbers that have been called during play
	PreCalled   map[int]bool   `json:"preCalled"` // All the numbers that have been called before play
	rand        *rand.Rand     `copier:"-"`
	KnownShapes map[int]bool
}

type Square struct {
	Number    int          `json:"number"`    // not really necessary, number is the key from Game->Square
	Called    bool         `json:"called"`    // FIXME necessary?
	PreCalled bool         `json:"preCalled"` // FIXME necessary?
	Needed    map[int]bool `json:"needed"`    // to handle multiple shapes
}

const FreeSquareIndex = 12

// why does this care about being an actual bingo board?
// there are 75 possible numbers to call
// the player has 24 of them plus the free square
// for free square, just ignore any shapes for it
// so... just set up a map for each shape? or a map of the numbers?

// Returns a Game with no squares, no calls, and an initialized RNG
func newGame() *Game {
	g := Game{}
	g.Squares = make(map[int]Square)
	g.Called = make(map[int]bool)
	g.PreCalled = make(map[int]bool)
	g.KnownShapes = make(map[int]bool)
	s1 := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(s1)
	return &g
}

func (g *Game) NumShapes() int {
	return len(g.KnownShapes)
}

// Add a square to the player's bingo board with its number and which shapes it's needed for
func (g *Game) addSquare(bingoNum int, preCalled bool, needed ...int) {
	sq := newSquare(bingoNum)
	sq.PreCalled = preCalled
	for _, n := range needed {
		sq.Needed[n] = true
		g.KnownShapes[n] = true
	}
	g.Squares[bingoNum] = sq
}

// Mark the square played, and check if the board has won.
func (g *Game) playSquare(bingoNum int) bool {

	square, exists := g.Squares[bingoNum]

	// See if we have the square for the number that was called
	if !exists { // we don't have that square
		return false
	}

	square.Called = true
	g.Squares[bingoNum] = square

	for shapeId, _ := range square.Needed {
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
		return 0, errors.New("no more numbers to pick")
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
	for _, square := range g.Squares {
		square.Called = square.PreCalled
		g.Squares[square.Number] = square
	}
}

func newSquare(n int) Square {
	return Square{Number: n, Called: false, PreCalled: false, Needed: make(map[int]bool)}
}
