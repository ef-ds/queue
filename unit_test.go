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

package queue

import (
	"testing"
)

const (
	refillCount = 3
	pushCount   = maxInternalSliceSize * 3 // Push to fill at least 3 internal slices
)

func TestNewShouldReturnInitiazedInstanceOfqueue(t *testing.T) {
	q := New()
	assertInvariants(t, q, nil)
}

func TestInvariantsWhenEmptyInMiddleOfSlice(t *testing.T) {
	q := new(Queue)
	q.Push(0)
	assertInvariants(t, q, nil)
	q.Push(1)
	assertInvariants(t, q, nil)
	q.Pop()
	assertInvariants(t, q, nil)
	q.Pop()
	// At this point, the queue is empty and hp will
	// not be pointing at the start of the slice.
	assertInvariants(t, q, nil)
}

func TestPushPopShouldHaveAllInternalLinksInARing(t *testing.T) {
	q := New()
	pushValue, extraAddedItems := 0, 0

	// Push maxInternalSliceSize items to fill the first array
	for i := 1; i <= maxInternalSliceSize; i++ {
		pushValue++
		q.Push(pushValue)
	}

	// Push 1 extra item to force the creation of a new array
	pushValue++
	q.Push(pushValue)
	extraAddedItems++
	checkLinks(t, q, pushValue, maxInternalSliceSize)

	// Push another maxInternalSliceSize-1 to fill the second array
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		pushValue++
		q.Push(pushValue)
		checkLinks(t, q, pushValue, maxInternalSliceSize)
	}

	// Push 1 extra item to force the creation of a new array (3 total)
	pushValue++
	q.Push(pushValue)
	checkLinks(t, q, pushValue, maxInternalSliceSize)

	// Check final len after all pushes
	if q.Len() != maxInternalSliceSize+maxInternalSliceSize+extraAddedItems {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize+maxInternalSliceSize+extraAddedItems, q.Len())
	}

	// Pop one item to force moving the tail to the middle slice. This also means the old tail
	// slice should have no items now
	expectedLen := q.Len()
	popValue := 1
	if v, ok := q.Pop(); !ok || v.(int) != popValue {
		t.Errorf("Expected: %d; Got: %d", popValue, v)
	}
	expectedLen--
	checkLinks(t, q, expectedLen, maxInternalSliceSize)

	// Pop maxInternalSliceSize-1 items to empty the tail (middle) slice
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		popValue++
		if v, ok := q.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		expectedLen--
		checkLinks(t, q, expectedLen, maxInternalSliceSize)
	}

	// Pop one extra item to force moving the tail to the head (first) slice. This also means the old tail
	// slice should have no items now.
	popValue++
	if v, ok := q.Pop(); !ok || v.(int) != popValue {
		t.Errorf("Expected: %d; Got: %d", popValue, v)
	}
	expectedLen--
	checkLinks(t, q, expectedLen, maxInternalSliceSize)

	// Pop maxFirstSliceSize-1 items to empty the head (first) slice
	for i := 1; i <= maxInternalSliceSize; i++ {
		popValue++
		if v, ok := q.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		expectedLen--
		checkLinks(t, q, expectedLen, maxInternalSliceSize)
	}

	// The queue shoud be empty
	if q.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, q.Len())
	}
	if _, ok := q.Front(); ok {
		t.Error("Expected: false; Got: true")
	}
	if cap(q.tail.v) != maxInternalSliceSize {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize, cap(q.tail.v))
	}
}

// Helper methods-----------------------------------------------------------------------------------

// Checks the internal slices and its linkq.
func checkLinks(t *testing.T, q *Queue, length, tailSliceSize int) {
	t.Helper()
	if q.Len() != length {
		t.Errorf("Unexpected length; Expected: %d; Got: %d", length, q.Len())
	}
	if cap(q.tail.v) != tailSliceSize {
		t.Errorf("Unexpected tail size; Expected: %d; Got: %d", tailSliceSize, len(q.tail.v))
	}
	if t.Failed() {
		t.FailNow()
	}
}

// assertInvariants checks all the invariant conditions in d that we can think of.
// If val is non-nil it is used to find the expected value for an item at index
// i measured from the head of the queue.
func assertInvariants(t *testing.T, q *Queue, val func(i int) interface{}) {
	t.Helper()
	fail := func(what string, got, want interface{}) {
		t.Errorf("invariant fail: %s; got %v want %v", what, got, want)
	}
	if q == nil {
		fail("non-nil queue", q, "non-nil")
	}
	if q.tail == nil {
		// Zero value.
		if q.tail != nil {
			fail("nil tail when zero", q.tail, nil)
		}
		if q.len != 0 {
			fail("zero length when zero", q.len, 0)
		}
		return
	}
	if t.Failed() {
		t.FailNow()
	}
}
