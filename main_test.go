package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_playUntilWin(t *testing.T) {

	var game = newGame()

	game.setNeeded(1, []int{1, 3, 5, 7, 9})
	game.setNeeded(2, []int{1, 2, 4, 6, 8, 10})

	calls, winners, err := playUntilWin(game)
	assert.Nil(t, err)
	assert.True(t, calls >= 5)
	assert.True(t, len(winners) >= 1)
}

func Test_playUntilWin_twoWinners(t *testing.T) {
	var game = newGame()

	game.setNeeded(1, []int{1, 2, 4, 6, 8, 10})
	game.setNeeded(2, []int{1, 2, 4, 6, 8, 10})
	calls, winners, err := playUntilWin(game)
	assert.Nil(t, err)
	assert.True(t, calls >= 6)
	assert.True(t, len(winners) == 2)
}

// This test is skipped because there are no callable squares and we only check wins after a call
func Test_freeSquare(t *testing.T) {
	//t.Skipf("No longer relevant")
	// if a shape only requires the free square it should win on one call
	var game = newGame()
	game.setNeeded(1, []int{12})
	calls, _, err := playUntilWin(game)
	assert.Nil(t, err)
	assert.True(t, calls == 0)
}

func Test_freeSquare_andOneMore(t *testing.T) {
	// if a shape only requires the free square and one more it should win instantly
	var game = newGame()
	game.setNeeded(1, []int{12, 14})
	won := game.playSquare(14)
	assert.True(t, won)
}

func Test_computeAveragePlaysUntilWin(t *testing.T) {
	var game = newGame()
	game.setNeeded(1, []int{1, 2, 4, 6, 8, 10})
	avgCalls, winningShapes := computeAveragePlaysUntilWin(game, 10)
	assert.True(t, avgCalls > 6)
	assert.True(t, avgCalls <= 75)
	assert.Equal(t, 10, winningShapes[1])
	log.Printf("avgCalls=%d, winningShapes=%v", avgCalls, winningShapes)
}

func Test_createGameFromDataMap(t *testing.T) {

	sqdata := []map[string]string{}
	sq1 := make(map[string]string)
	sq1["name"] = "square_id_1"
	sq1["value"] = "42"
	sqdata = append(sqdata, sq1)
	sq1s := make(map[string]string)
	sq1s["name"] = "square_needed_1_2"
	sq1s["value"] = "on"
	sqdata = append(sqdata, sq1s)
	sq1c := make(map[string]string)
	sq1c["name"] = "called_42"
	sq1c["value"] = "on"
	sqdata = append(sqdata, sq1s)

	g, err := createGameFromDataMap(sqdata)
	assert.Nil(t, err)
	assert.Equal(t, 42, g.Squares[1].Number)
	assert.False(t, g.Squares[1].Needed[1])
	assert.True(t, g.Squares[1].Needed[2])
	assert.Equal(t, 42, g.Squares[1].Number)

}
