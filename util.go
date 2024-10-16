package main

type int_tuple struct {
	first  int
	second int
}

func remove(s []tile, i int) []tile {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
