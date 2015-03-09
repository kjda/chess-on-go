package main

type Square uint8

//get rank from square
func (s Square) rank() int {
	return int(s / 8)
}

//get file from square
func (s Square) file() int {
	return int(s % 8)
}
