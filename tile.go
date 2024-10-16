package main

import (
	"fmt"
)

type tile_side int

const (
	grass tile_side = iota
	road
	city
	city_connected
)

const (
	left int = iota
	top
	right
	bottom
)

type tile struct {
	Crest     bool
	Ends_road bool
	Cathedral bool
	Meeple    bool
	Sides     [4]tile_side
	Neighbors [6]*tile
}

func sides_match(first tile_side, second tile_side) bool {
	are_same := first == second
	both_city := (first == city && second == city_connected) || (first == city_connected && second == city)
	if are_same || both_city {
		return true
	}

	return false
}

func make_tile(left tile_side, top tile_side, right tile_side, bottom tile_side, crest bool, ends_road bool, cahedral bool) tile {
	return tile{
		Crest:     crest,
		Ends_road: ends_road,
		Cathedral: cahedral,
		Sides:     [4]tile_side{left, top, right, bottom},
	}
}

func make_tile_short(left tile_side, top tile_side, right tile_side, bottom tile_side) tile {
	return tile{
		Crest:     false,
		Ends_road: false,
		Cathedral: false,
		Sides:     [4]tile_side{left, top, right, bottom},
	}
}

func make_tile_from_array(sides [4]tile_side, crest bool, ends_road bool, cathedral bool) tile {
	return tile{
		Crest:     crest,
		Ends_road: ends_road,
		Cathedral: cathedral,
		Sides:     sides,
	}
}

func (t tile) print() {
	fmt.Printf("%+v\n", t)
}

func (t *tile) rotate() {
	t.Sides = [4]tile_side{t.Sides[3], t.Sides[0], t.Sides[1], t.Sides[2]}
}
