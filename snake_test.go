package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnake(t *testing.T) {
	m := snake{
		points: []point{
			{10, 10},
		},
	}
	m.grow(right)
	assert.Equal(t, []point{
		{11, 10}, {10, 10},
	}, m.getAll())

	m.grow(right)
	assert.Equal(t, []point{
		{12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.grow(left)
	assert.Equal(t, []point{
		{12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.grow(down)
	assert.Equal(t, []point{
		{12, 11}, {12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.grow(right)
	assert.Equal(t, []point{
		{13, 11}, {12, 11}, {12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.grow(up)
	assert.Equal(t, []point{
		{13, 10}, {13, 11}, {12, 11}, {12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.grow(up)
	assert.Equal(t, []point{
		{13, 9}, {13, 10}, {13, 11}, {12, 11}, {12, 10}, {11, 10}, {10, 10},
	}, m.getAll())

	m.move(up)
	assert.Equal(t, []point{
		{13, 8}, {13, 9}, {13, 10}, {13, 11}, {12, 11}, {12, 10}, {11, 10},
	}, m.getAll())

	m.move(down)
	assert.Equal(t, []point{
		{13, 8}, {13, 9}, {13, 10}, {13, 11}, {12, 11}, {12, 10}, {11, 10},
	}, m.getAll())

	m.move(left)
	assert.Equal(t, []point{
		{12, 8}, {13, 8}, {13, 9}, {13, 10}, {13, 11}, {12, 11}, {12, 10},
	}, m.getAll())

}
