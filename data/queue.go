package data

import "github.com/soupstoregames/gamelib/utils"

// Queue is your standard First In First Out queue.
type Queue[T any] struct {
	data []T
}

func (q *Queue[T]) Push(e T) {
	q.data = append(q.data, e)
}

func (q *Queue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		return utils.Zero[T](), false
	}
	e := q.data[0]
	q.data = q.data[1:]
	return e, true
}

func (q *Queue[T]) Peek(i int) T {
	return q.data[i]
}

func (q *Queue[T]) Len() int {
	return len(q.data)
}
