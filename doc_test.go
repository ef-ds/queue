package queue_test

import (
	"fmt"

	"github.com/ef-ds/queue"
)

func Example() {
	var q queue.Queue

	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	for q.Len() > 0 {
		v, _ := q.Pop()
		fmt.Print(v)
	}
	// Output: 12345
}
