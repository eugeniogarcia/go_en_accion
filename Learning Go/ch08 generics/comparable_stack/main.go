package main

import (
	"fmt"
)

// restringimos T a tipos comparables usando comparable
type Stack[T comparable] struct {
	vals []T
}

func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T //para retornar un zero del tipo T, declarando una variable de tipo T sin inicializar y la devolvemos
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func (s Stack[T]) Contains(val T) bool {
	for _, v := range s.vals {
		if v == val { //podemos usar == porque T es comparable
			return true
		}
	}
	return false
}

func main() {
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	s.Push(30)
	fmt.Println(s.Contains(10))
	fmt.Println(s.Contains(5))
}
