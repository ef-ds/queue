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

// Package queue implements a very fast and efficient general purpose
// First-In-First-Out (FIFO) queue data structure that is specifically optimized
// to perform when used by Microservices and serverless services running in
// production environments.
package queue

const (
	// firstSliceSize holds the size of the first slice.
	firstSliceSize = 4

	// sliceGrowthFactor determines by how much and how fast the first internal
	// slice should grow. A growth factor of 4, firstSliceSize = 4 and maxFirstSliceSize = 64,
	// the first slice will start with size 4, then 16 (4*4), then 64 (16*4).
	// The growth factor should be tweaked together with firstSliceSize and specially,
	// maxFirstSliceSize for maximum efficiency.
	// sliceGrowthFactor only applies to the very first slice created. All other
	// subsequent slices are created with fixed size of maxInternalSliceSize.
	sliceGrowthFactor = 4

	// maxFirstSliceSize holds the maximum size of the first slice.
	maxFirstSliceSize = 64

	// maxInternalSliceSize holds the maximum size of each internal slice.
	maxInternalSliceSize = 256
)

// Queue implements an unbounded, dynamically growing double-ended-queue (queue).
// The zero value for queue is an empty queue ready to use.
type Queue struct {
	// Head points to the first node of the linked list.
	head *node

	// Tail points to the last node of the linked list.
	// In an empty queue, head and tail points to the same node.
	tail *node

	// Hp is the index pointing to the current first element in the queue
	// (i.e. first element added in the current queue values).
	hp int

	// hlp points to the last index in the head slice.
	hlp int

	// tp is the index pointing one beyond the current last element in the queue
	// (i.e. last element added in the current queue values).
	tp int

	// Len holds the current queue values length.
	len int
}

// Node represents a queue node.
// Each node holds a slice of user managed values.
type node struct {
	// v holds the list of user added values in this node.
	v []interface{}

	// n points to the next node in the linked list.
	n *node
}

// New returns an initialized queue.
func New() *Queue {
	return new(Queue)
}

// Init initializes or clears queue d.
func (d *Queue) Init() *Queue {
	*d = Queue{}
	return d
}

// Len returns the number of elements of queue d.
// The complexity is O(1).
func (d *Queue) Len() int { return d.len }

// Front returns the first element of queue d or nil if the queue is empty.
// The second, bool result indicates whether a valid value was returned;
// if the queue is empty, false will be returned.
// The complexity is O(1).
func (d *Queue) Front() (interface{}, bool) {
	if d.len == 0 {
		return nil, false
	}
	return d.head.v[d.hp], true
}

// Push adds value v to the the back of the queue.
// The complexity is O(1).
func (d *Queue) Push(v interface{}) {
	switch {
	case d.head == nil:
		// No nodes present yet.
		h := &node{v: make([]interface{}, firstSliceSize)}
		h.n = h
		d.head = h
		d.tail = h
		d.tail.v[0] = v
		d.hlp = firstSliceSize - 1
		d.tp = 1
	case d.tp < len(d.tail.v):
		// There's room in the tail slice.
		d.tail.v[d.tp] = v
		d.tp++
	case d.tp < maxFirstSliceSize:
		// We're on the first slice and it hasn't grown large enough yet.
		nv := make([]interface{}, len(d.tail.v)*sliceGrowthFactor)
		copy(nv, d.tail.v)
		d.tail.v = nv
		d.tail.v[d.tp] = v
		d.tp++
		d.hlp = len(nv) - 1
	case d.tail.n != d.head:
		// There's at least one spare link between head and tail nodes.
		n := d.tail.n
		d.tail = n
		d.tail.v[0] = v
		d.tp = 1
	default:
		// No available nodes, so make one.
		n := &node{v: make([]interface{}, maxInternalSliceSize)}
		n.n = d.head
		d.tail.n = n
		d.tail = n
		d.tail.v[0] = v
		d.tp = 1
	}
	d.len++
}

// Pop retrieves and removes the current element from the front of the queue.
// The second, bool result indicates whether a valid value was returned;
// if the queue is empty, false will be returned.
// The complexity is O(1).
func (d *Queue) Pop() (interface{}, bool) {
	if d.len == 0 {
		return nil, false
	}

	vp := &d.head.v[d.hp]
	v := *vp
	*vp = nil // Avoid memory leaks
	d.len--
	switch {
	case d.hp < d.hlp:
		// The head isn't at the end of the slice, so just
		// move on one place.
		d.hp++
	case d.head == d.tail:
		// There's only a single element at the end of the slice
		// so we can't increment hp, so change tp instead.
		d.tp = d.hp
	default:
		// Move to the next slice.
		d.hp = 0
		d.head = d.head.n
		d.hlp = len(d.head.v) - 1
	}
	return v, true
}
