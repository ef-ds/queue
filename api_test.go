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

func TestPopWithZeroValueShouldReturnReadyToUsequeue(t *testing.T) {
	var q queue.Queue
	q.Push(1)
	q.Push(2)

	v, ok := q.Front()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	v, ok = q.Pop()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	v, ok = q.Front()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	v, ok = q.Pop()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	_, ok = q.Front()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
	_, ok = q.Pop()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
}

func TestWithZeroValueAndEmptyShouldReturnAsEmpty(t *testing.T) {
	var q queue.Queue

	if _, ok := q.Front(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := q.Front(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := q.Pop(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if l := q.Len(); l != 0 {
		t.Errorf("Expected: 0 as the queue is empty; Got: %d", l)
	}
}

func TestInitShouldReturnEmptyqueue(t *testing.T) {
	var q queue.Queue
	q.Push(1)

	q.Init()

	if _, ok := q.Front(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := q.Pop(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if l := q.Len(); l != 0 {
		t.Errorf("Expected: 0 as the queue is empty; Got: %d", l)
	}
}

func TestPopWithNilValuesShouldReturnAllValuesInOrder(t *testing.T) {
	q := queue.New()
	q.Push(1)
	q.Push(nil)
	q.Push(2)
	q.Push(nil)

	v, ok := q.Pop()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	v, ok = q.Pop()
	if !ok || v != nil {
		t.Errorf("Expected: nil; Got: %d", v)
	}
	v, ok = q.Pop()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	v, ok = q.Pop()
	if !ok || v != nil {
		t.Errorf("Expected: nil; Got: %d", v)
	}
	_, ok = q.Pop()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
}
