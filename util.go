package main

type int_tuple struct {
	first  int
	second int
}

func addTuples(x int_tuple, y int_tuple) int_tuple {
	return int_tuple{x.first + y.first, x.second + y.second}
}

func remove(s []tile, i int) []tile {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
