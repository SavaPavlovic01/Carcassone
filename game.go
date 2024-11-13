package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var gggg [4]tile_side = [4]tile_side{grass, grass, grass, grass}
var gggr [4]tile_side = [4]tile_side{grass, grass, grass, road}
var oooo [4]tile_side = [4]tile_side{city_connected, city_connected, city_connected, city_connected}
var ooog [4]tile_side = [4]tile_side{city_connected, city_connected, city_connected, grass}
var ooor [4]tile_side = [4]tile_side{city_connected, city_connected, city_connected, road}
var oogg [4]tile_side = [4]tile_side{city_connected, city_connected, grass, grass}
var oorr [4]tile_side = [4]tile_side{city_connected, city_connected, road, road}
var ogog [4]tile_side = [4]tile_side{city_connected, grass, city_connected, grass}
var ccgg [4]tile_side = [4]tile_side{city, city, grass, grass}
var gcgc [4]tile_side = [4]tile_side{grass, city, grass, city}
var gcgg [4]tile_side = [4]tile_side{grass, city, grass, grass}
var rcgr [4]tile_side = [4]tile_side{road, city, grass, road}
var gcrr [4]tile_side = [4]tile_side{grass, city, road, road}
var rcrr [4]tile_side = [4]tile_side{road, city, road, road}
var rcrg [4]tile_side = [4]tile_side{road, city, road, grass}
var grgr [4]tile_side = [4]tile_side{grass, road, grass, road}
var rggr [4]tile_side = [4]tile_side{road, grass, grass, road}
var rgrr [4]tile_side = [4]tile_side{road, grass, road, road}
var rrrr [4]tile_side = [4]tile_side{road, road, road, road}

type Meeple struct {
	x        int
	y        int
	color    string
	isPriest bool
}
type game struct {
	Deck    []tile
	Board   map[int_tuple]tile
	Meeples []Meeple
}

// TODO: REFACTOR
func init_deck() []tile {
	tiles := []tile{}

	for i := 0; i < 4; i++ {
		tiles = append(tiles, make_tile_from_array(gggg, false, true, true))
		tiles = append(tiles, make_tile_from_array(rgrr, false, true, false))
	}

	for i := 0; i < 3; i++ {
		tiles = append(tiles, make_tile_from_array(ooog, false, false, false))
		tiles = append(tiles, make_tile_from_array(oogg, false, false, false))
		tiles = append(tiles, make_tile_from_array(oorr, false, false, false))
		tiles = append(tiles, make_tile_from_array(gcgc, false, false, false))
		tiles = append(tiles, make_tile_from_array(rcgr, false, false, false))
		tiles = append(tiles, make_tile_from_array(gcrr, false, false, false))
		tiles = append(tiles, make_tile_from_array(rcrr, false, false, false))
		tiles = append(tiles, make_tile_from_array(rcrg, false, false, false))
	}

	for i := 0; i < 2; i++ {
		tiles = append(tiles, make_tile_from_array(gggr, false, true, true))
		tiles = append(tiles, make_tile_from_array(ooor, true, true, false))
		tiles = append(tiles, make_tile_from_array(oogg, true, false, false))
		tiles = append(tiles, make_tile_from_array(oorr, true, false, false))
		tiles = append(tiles, make_tile_from_array(ogog, true, false, false))
		tiles = append(tiles, make_tile_from_array(ccgg, false, false, false))
	}

	tiles = append(tiles, make_tile_from_array(oooo, true, false, false))
	tiles = append(tiles, make_tile_from_array(ooog, true, false, false))
	tiles = append(tiles, make_tile_from_array(ooor, false, true, false))
	tiles = append(tiles, make_tile_from_array(ogog, false, false, false))
	tiles = append(tiles, make_tile_from_array(rrrr, false, true, false))

	for range 5 {
		tiles = append(tiles, make_tile_from_array(gcgg, false, false, false))
		tiles = append(tiles, make_tile_from_array(grgr, false, false, false))
		tiles = append(tiles, make_tile_from_array(rggr, false, false, false))
	}

	for range 3 {
		tiles = append(tiles, make_tile_from_array(grgr, false, false, false))
		tiles = append(tiles, make_tile_from_array(rggr, false, false, false))
	}

	tiles = append(tiles, make_tile_from_array(rggr, false, false, false))

	return tiles
}

func new_game() game {
	return game{
		Deck:  init_deck(),
		Board: map[int_tuple]tile{{0, 0}: make_tile_from_array(rcrg, false, false, false)},
	}
}

func (g *game) draw_card() (drawn_tile tile) {
	index := rand.Intn(len(g.Deck))
	drawn_tile = g.Deck[index]
	g.Deck = remove(g.Deck, index)
	return drawn_tile
}

func (g game) place_tile(placed_tile tile, x_offset int, y_offset int) error {
	_, ok := g.Board[int_tuple{x_offset, y_offset}]
	if ok {
		return errors.New("location occupied")
	}

	if !g.check_neighbors(placed_tile, x_offset, y_offset) {
		return errors.New("invalid location")
	}

	g.Board[int_tuple{x_offset, y_offset}] = placed_tile
	return nil
}

func getOppositeSide(side int) int {
	return (side + 2) % 4
}

func (g game) check_neighbors(t tile, x_offset int, y_offset int) bool {
	//fmt.Println(x_offset, y_offset)
	at_least_one_neighbor := false

	offset := []int_tuple{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	sides := []int{left, top, right, bottom}
	location := int_tuple{x_offset, y_offset}
	fmt.Println(t)
	for i, side := range sides {
		neighbor, exists := g.Board[addTuples(location, offset[i])]
		fmt.Println(side, neighbor)
		if !exists {
			continue
		}
		at_least_one_neighbor = true
		if !sides_match(t.Sides[side], neighbor.Sides[getOppositeSide(side)]) {
			return false
		}
	}

	return at_least_one_neighbor
}

func (g *game) addMeeple(x int, y int, color string, isPriest bool) {
	g.Meeples = append(g.Meeples, Meeple{x: x, y: y, color: color, isPriest: isPriest})
}

func (g *game) removeMeeple(pos int) {
	g.Meeples = append(g.Meeples[:pos], g.Meeples[pos+1:]...)
}

func (g *game) moveMeeple(index int, x int, y int) {
	g.Meeples[index].x = x
	g.Meeples[index].y = y
}
