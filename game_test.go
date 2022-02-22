package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame_winner(t *testing.T) {

	var game = newGame()
	for i := 1; i <= 10; i++ {
		game.addSquare(i, false, i%2)
	}

	assert.False(t, game.playSquare(1))
	assert.False(t, game.playSquare(2))
	assert.False(t, game.playSquare(3))
	assert.False(t, game.playSquare(5))
	assert.False(t, game.playSquare(7))
	assert.True(t, game.playSquare(9))
	assert.False(t, game.winner(0))
	assert.True(t, game.winner(1))

}

func TestGame_pickUncalledSquare(t *testing.T) {
	var game = newGame()
	n, err := game.callRandomSquare()
	assert.True(t, n >= 1)
	assert.True(t, n <= 75)
	assert.Nil(t, err)
	// Call the rest of the board
	for i := 1; i <= 74; i++ {
		m, err := game.callRandomSquare()
		assert.False(t, n == m)
		assert.True(t, m >= 1)
		assert.True(t, m <= 75)
		assert.Nil(t, err)
	}
	n, err = game.callRandomSquare()
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

}
