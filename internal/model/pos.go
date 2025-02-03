package model

type Pos [2]int

func (p Pos) I() int {
	return p[0]
}

func (p Pos) J() int {
	return p[1]
}
