// Copyright (c) 2018 ef-ds
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package queue_test

import (
	"testing"

	"github.com/ef-ds/queue"
)

const (
	pushCount = 256 * 3 // Push to fill at least 3 internal slices
)

func TestFillQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue

	for i := 0; i < pushCount; i++ {
		q.Push(i)
	}
	for i := 0; i < pushCount; i++ {
		if v, ok := q.Pop(); !ok || v.(int) != i {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, q.Len())
	}
}

func TestRefillQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue

	for i := 0; i < refillCount; i++ {
		for j := 0; j < pushCount; j++ {
			q.Push(j)
		}
		for j := 0; j < pushCount; j++ {
			if v, ok := q.Pop(); !ok || v.(int) != j {
				t.Errorf("Expected: %d; Got: %d", i, v)
			}
		}
		if q.Len() != 0 {
			t.Errorf("Expected: %d; Got: %d", 0, q.Len())
		}
	}
}

func TestRefillFullQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue
	for i := 0; i < pushCount; i++ {
		q.Push(i)
	}

	for i := 0; i < refillCount; i++ {
		for j := 0; j < pushCount; j++ {
			q.Push(j)
		}
		for j := 0; j < pushCount; j++ {
			if v, ok := q.Pop(); !ok || v.(int) != j {
				t.Errorf("Expected: %d; Got: %d", j, v)
			}
		}
		if q.Len() != pushCount {
			t.Errorf("Expected: %d; Got: %d", pushCount, q.Len())
		}
	}
}

func TestSlowIncreaseQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue

	count := 0
	for i := 0; i < pushCount; i++ {
		count++
		q.Push(count)
		count++
		q.Push(count)
		if v, ok := q.Pop(); !ok || v.(int) != i+1 {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if q.Len() != pushCount {
		t.Errorf("Expected: %d; Got: %d", pushCount, q.Len())
	}
}

func TestSlowDecreaseQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue
	push := 0
	for i := 0; i < pushCount; i++ {
		q.Push(push)
		push++
	}

	count := -1
	for i := 0; i < pushCount-1; i++ {
		count++
		if v, ok := q.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}
		count++
		if v, ok := q.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}

		q.Push(push)
		push++
	}
	count++
	if v, ok := q.Pop(); !ok || v.(int) != count {
		t.Errorf("Expected: %d; Got: %d", count, v)
	}
	if q.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, q.Len())
	}
}

func TestStableQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue

	for i := 0; i < pushCount; i++ {
		q.Push(i)
		if v, ok := q.Pop(); !ok || v.(int) != i {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, q.Len())
	}
}

func TestStableFullQueueShouldRetrieveAllElementsInOrder(t *testing.T) {
	var q queue.Queue

	for i := 0; i < pushCount; i++ {
		q.Push(i)
		if v, ok := q.Pop(); !ok || v.(int) != i {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, q.Len())
	}
}

func TestPushFrontPopRefillWith0ToPushCountItemsShouldReturnAllValuesInOrder(t *testing.T) {
	var q queue.Queue

	for i := 0; i < refillCount; i++ {
		for k := 0; k < pushCount; k++ {
			for j := 0; j < k; j++ {
				q.Push(j)
			}
			for j := 0; j < k; j++ {
				v, ok := q.Pop()
				if !ok || v == nil || v.(int) != j {
					t.Errorf("Expected: %d; Got: %d", j, v)
				}
			}
			if q.Len() != 0 {
				t.Errorf("Expected: %d; Got: %d", 0, q.Len())
			}
		}
	}
}
