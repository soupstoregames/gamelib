package data

import "github.com/soupstoregames/gamelib/utils"

// Stack is your standard First In Last Out stack.
type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(e T) {
	s.data = append(s.data, e)
}

func (s *Stack[T]) Pop() T {
	if len(s.data) == 0 {
		return utils.Zero[T]()
	}
	end := len(s.data) - 1
	e := s.data[end]
	s.data = s.data[:end]
	return e
}

func (s *Stack[T]) Peek() T {
	if len(s.data) == 0 {
		return utils.Zero[T]()
	}
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}
