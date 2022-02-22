package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// creates a game with 10 squares and 2 shapes; one shape is all evens, one all odds
func createSimpleGame() *Game {
	var game = newGame()
	for i := 1; i <= 10; i++ {
		game.addSquare(i, false, i%2)
	}
	return game
}

func Test_playUntilWin(t *testing.T) {

	game := createSimpleGame()
	calls, winners, err := playUntilWin(game)
	assert.Nil(t, err)
	assert.True(t, calls >= 5)
	assert.True(t, len(winners) >= 1)
}

func Test_playUntilWin_twoWinners(t *testing.T) {
	// game has 2 shapes, both all evens
	var game = newGame()
	for i := 1; i <= 10; i++ {
		game.addSquare(i, false, 0, 1)
	}

	calls, winners, err := playUntilWin(game)
	assert.Nil(t, err)
	assert.True(t, calls >= 5)
	assert.True(t, len(winners) == 2)
}

func Test_computeAveragePlaysUntilWin(t *testing.T) {
	// Game with only one shape
	var game = newGame()
	for i := 1; i <= 10; i++ {
		game.addSquare(i, false, i%2)
	}
	// and a fake shape that can't win, because we must have 2 shapes
	sq := newSquare(100)
	sq.Needed[1] = true
	game.Squares[100] = sq

	avgCalls, winningShapes := computeAveragePlaysUntilWin(game, 10)
	assert.True(t, avgCalls > 6)
	assert.True(t, avgCalls <= 75)
	assert.Equal(t, 10, winningShapes[0])
	log.Printf("avgCalls=%d, winningShapes=%v", avgCalls, winningShapes)
}

func Test_createGameFromDataMap(t *testing.T) {

	sqdata := []map[string]string{}
	sq1 := make(map[string]string)
	sq1["name"] = "square_id_1"
	sq1["value"] = "42"
	sqdata = append(sqdata, sq1)
	sq1s := make(map[string]string)
	sq1s["name"] = "square_needed_1_1"
	sq1s["value"] = "on"
	sqdata = append(sqdata, sq1s)
	sq1c := make(map[string]string)
	sq1c["name"] = "called_42"
	sq1c["value"] = "on"
	sqdata = append(sqdata, sq1s)

	g, err := createGameFromDataMap(sqdata)
	assert.Nil(t, err)
	assert.Equal(t, 42, g.Squares[42].Number)
	assert.False(t, g.Squares[42].Needed[0])
	assert.True(t, g.Squares[42].Needed[1])

}
