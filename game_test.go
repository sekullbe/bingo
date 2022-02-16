package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame_winner(t *testing.T) {

	var game = newGame()
	game.setNeeded(1, []int{1, 3, 5, 7, 9})
	game.setNeeded(2, []int{2, 4, 6, 8, 10})
	assert.False(t, game.playSquare(1))
	assert.False(t, game.playSquare(2))
	assert.False(t, game.playSquare(3))
	assert.False(t, game.playSquare(5))
	assert.False(t, game.playSquare(7))
	assert.True(t, game.playSquare(9))
	assert.True(t, game.winner(1))
	assert.False(t, game.winner(2))
}

func TestGame_pickUncalledSquare(t *testing.T) {
	var game = newGame()
	n, err := game.callRandomSquare()
	assert.True(t, n >= 1)
	assert.True(t, n <= 75)
	assert.Nil(t, err)
	for i := 1; i <= 74; i++ { // pick 74 more numbers
		m, err := game.callRandomSquare()
		assert.False(t, n == m)
		assert.True(t, m >= 1)
		assert.True(t, m <= 75)
		assert.Nil(t, err)
	}
	n, err = game.callRandomSquare()
	assert.NotNil(t, err)
	assert.Equal(t, 0, n)

}
