# queue [![Build Status](https://travis-ci.com/ef-ds/queue.svg?branch=master)](https://travis-ci.com/ef-ds/queue) [![codecov](https://codecov.io/gh/ef-ds/queue/branch/master/graph/badge.svg)](https://codecov.io/gh/ef-ds/queue) [![Go Report Card](https://goreportcard.com/badge/github.com/ef-ds/queue)](https://goreportcard.com/report/github.com/ef-ds/queue)  [![GoDoc](https://godoc.org/github.com/ef-ds/queue?status.svg)](https://godoc.org/github.com/ef-ds/queue)

Package queue implements a very fast and efficient general purpose First-In-First-Out (FIFO) queue data structure that is specifically optimized to perform when used by Microservices and serverless services running in production environments. Internally, queue stores the elements in a dynamic growing circular singly linked list of arrays.


## Install
From a configured [Go environment](https://golang.org/doc/install#testing):
```sh
go get -u github.com/ef-ds/queue
```

If you are using dep:
```sh
dep ensure -add github.com/ef-ds/queue@1.0.0
```

We recommend to target only released versions for production use.


## How to Use
```go
package main

import (
	"fmt"

	"github.com/ef-ds/queue"
)

func main() {
	var q queue.Queue

	for i := 1; i <= 5; i++ {
		q.Push(i)
	}
	for q.Len() > 0 {
		v, _ := q.Pop()
		fmt.Println(v)
	}
}
```

Output:
```
1
2
3
4
5
```

Also refer to the [integration](integration_test.go) and [API](api_test.go) tests.



## Tests
Besides having 100% code coverage, queue has an extensive set of [unit](unit_test.go), [integration](integration_test.go) and [API](api_test.go) tests covering all happy, sad and edge cases.

When considering all tests, queue has over 4x more lines of testing code when compared to the actual, functional code.

Performance and efficiency are major concerns, so queue has an extensive set of benchmark tests as well comparing the queue performance with a variety of high quality open source queue implementations.

See the [benchmark tests](https://github.com/ef-ds/queue-bench-tests/blob/master/BENCHMARK_TESTS.md) for details.


## Performance
Queue has constant time (O(1)) on all its operations (Push/Pop/Back/Len). It's not amortized constant because queue never copies more than 64 (maxInternalSliceSize/sliceGrowthFactor) items and when it expands or grow, it never does so by more than 256 (maxInternalSliceSize) items in a single operation.

Queue offers either the best or very competitive performance across all test sets, suites and ranges.

As a general purpose FIFO queue, queue offers, by far, the most balanced and consistent performance of all tested data structures.

See [performance](https://github.com/ef-ds/queue-bench-tests/blob/master/PERFORMANCE.md) for details.


## Design
The Efficient Data Structures (ef-ds) queue employs a new, modern queue design: a dynamic growing circular singly linked list of arrays.

That means the [FIFO queue](https://en.wikipedia.org/wiki/Queue_(abstract_data_type)) is a [singly-linked list](https://en.wikipedia.org/wiki/Singly_linked_list) where each node value is a fixed size [slice](https://tour.golang.org/moretypes/7). It is ring in shape because the linked list is a [circular one](https://en.wikipedia.org/wiki/Circular_buffer), where the last node always points to the first one in the ring.

![ns/op](testdata/queue.jpg?raw=true "queue Design")


### Design Considerations
Queue uses linked slices as its underlying data structure. The reason for the choice comes from two main observations of slice based queues:

1. When the queue needs to expand to accommodate new values, [a new, larger slice needs to be allocated](https://en.wikipedia.org/wiki/Dynamic_array#Geometric_expansion_and_amortized_cost) and used
2. Allocating and managing large slices is expensive, especially in an overloaded system with little available physical memory

To help clarify the scenario, below is what happens when a slice based queue that already holds, say 1bi items, needs to expand to accommodate a new item.

Slice based implementation.

- Allocate a new, twice the size of the previous allocated one, say 2 billion positions slice
- Copy over all 1 billion items from the previous slice into the new one
- Add the new value into the first unused position in the new slice, position 1000000001

The same scenario for queue plays out like below.

- Allocate a new 256 size slice
- Set the previous and next pointers
- Add the value into the first position of the new slice, position 0

The decision to use linked slices was also the result of the observation that slices goes to great length to provide predictive, indexed positions. A hash table, for instance, absolutely need this property, but not a queue. So queue completely gives up this property and focus on what really matters: add and retrieve from the edges (front/back). No copying around and repositioning of elements is needed for that. So when a slice goes to great length to provide that functionality, the whole work of allocating new arrays, copying data around is all wasted work. None of that is necessary. And this work costs dearly for large data sets as observed in the tests.

While linked slices design solve the slice expansion problem very effectively, it doesn't help with many real world usage scenarios such as in a stable processing environment where small amount of items are pushed and popped from the queue in a sequential way. This is a very common scenario for [Microservices](https://en.wikipedia.org/wiki/Microservices) and [serverless](https://en.wikipedia.org/wiki/Serverless_computing) services, for instance, where the service is able to handle the current traffic without stress.

To address the stable scenario in an effective way, queue keeps its internal linked arrays in a circular, ring shape. This way when items are pushed to the queue after some of them have been removed, the queue will automatically move over its tail slice back to the old head of the queue, effectively reusing the same already allocated slice. The result is a queue that will run through its ring reusing the ring to store the new values, instead of allocating new slices for the new values.



## Supported Data Types
Similarly to Go's standard library list, [list](https://github.com/golang/go/tree/master/src/container/list),
[ring](https://github.com/golang/go/tree/master/src/container/ring) and [heap](https://github.com/golang/go/blob/master/src/container/heap/heap.go) packages, queue supports "interface{}" as its data type. This means it can be used with any Go data types, including int, float, string and any user defined structs and pointers to interfaces.

The data types pushed into the queue can even be mixed, meaning, it's possible to push ints, floats and struct instances into the same queue.


## Safe for Concurrent Use
Queue is not safe for concurrent use. However, it's very easy to build a safe for concurrent use version of the queue. Impl7 design document includes an example of how to make impl7 safe for concurrent use using a mutex. queue can be made safe for concurret use using the same technique. Impl7 design document can be found [here](https://github.com/golang/proposal/blob/master/design/27935-unbounded-queue-package.md).


## Range Support
Just like the current container data structures such as [list](https://github.com/golang/go/tree/master/src/container/list),
[ring](https://github.com/golang/go/tree/master/src/container/ring) and [heap](https://github.com/golang/go/blob/master/src/container/heap/heap.go), queue doesn't support the range keyword for navigation.

However, the API offers two ways to iterate over the queue items. Either use "PopFront"/"PopBack" to retrieve the first current element and the second bool parameter to check for an empty queue.

```go
for v, ok := s.Pop(); ok; v, ok = s.Pop() {
    // Do something with v
}
```

Or use "Len" and "Pop" to check for an empty queue and retrieve the first current element.
```go
for s.Len() > 0 {
    v, _ := s.Pop()
    // Do something with v
}
```



## Why
We feel like this world needs improving. Our goal is to change the world, for the better, for everyone.

As software engineers at ef-ds, we feel like the best way we can contribute to a better world is to build amazing systems,
systems that solve real world problems, with unheard performance and efficiency.

We believe in challenging the status-quo. We believe in thinking differently. We believe in progress.

What if we could build queues, queues, lists, arrays, hash tables, etc that are much faster than the current ones we have? What if we had a dynamic array data structure that offers near constant time deletion (anywhere in the array)? Or that could handle 1 million items data sets using only 1/3 of the memory when compared to all known current implementations? And still runs 2x as fast?

One sofware engineer can't change the world him/herself, but a whole bunch of us can! Please join us improving this world. All the work done here is made 100% transparent and is 100% free. No strings attached. We only require one thing in return: please consider benefiting from it; and if you do so, please let others know about it.


## Competition
We're extremely interested in improving queue and we're on an endless quest for better efficiency and more performance. Please let us know your suggestions for possible improvements and if you know of other high performance queues not tested here, let us know and we're very glad to benchmark them.


## Releases
We're committed to a CI/CD lifecycle releasing frequent, but only stable, production ready versions with all proper tests in place.

We strive as much as possible to keep backwards compatibility with previous versions, so breaking changes are a no-go.

For a list of changes in each released version, see [CHANGELOG.md](CHANGELOG.md).


## Supported Go Versions
See [supported_go_versions.md](https://github.com/ef-ds/docs/blob/master/supported_go_versions.md).


## License
MIT, see [LICENSE](LICENSE).

"Use, abuse, have fun and contribute back!"


## Contributions
See [CONTRIBUTING.md](CONTRIBUTING.md).


## Roadmap
- Build tool to help find out the combination of firstSliceSize, sliceGrowthFactor, maxFirstSliceSize and maxInternalSliceSize that will yield the best performance
- Find the fastest open source queues and add them the bench tests
- Improve queue performance and/or efficiency by improving its design and/or implementation
- Build a high performance safe for concurrent use version of queue


## Contact
Suggestions, bugs, new queues to benchmark, issues with the queue, please let us know at ef-ds@outlook.com.
