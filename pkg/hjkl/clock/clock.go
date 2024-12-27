// Package clock contains an implementation of a delta clock.
package clock

// node stores events for a particular delta in a Clock.
type node[T comparable] struct {
	delta  int
	link   *node[T]
	events map[T]struct{}
}

// Clock implements a delta clock data structure, which is a linked list in
// which each node stores both a collection of events and a delta. The deltas
// are relative to the previous nodes in the list, so updating the entire clock
// an be done in O(1) time by decrementing or removing the head node delta.
// Adding new events can be done in O(n) time, where n is the number of nodes
// (not the number of events).
type Clock[T comparable] struct {
	head  *node[T]
	nodes map[T]*node[T]
}

// New creates an empty Clock.
func New[T comparable]() *Clock[T] {
	return &Clock[T]{nil, make(map[T]*node[T])}
}

// Schedule adds an event to the queue at the given delta.
func (c *Clock[T]) Schedule(t T, delta int) {
	c.Unschedule(t)

	var prev, curr, next *node[T] = nil, nil, c.head

	// Find where to insert event by iterating over the linked list of nodes,
	// ensuring we don't go past the end or past the desired node.
	for next != nil && delta > next.delta {
		delta -= next.delta
		prev, next = next, next.link
	}

	if next != nil && delta == next.delta {
		// If the desired node already exists, just reuse it.
		curr = next
	} else {
		// Desired node didn't exist, so create it with link to next node.
		curr = &node[T]{delta, next, make(map[T]struct{})}

		if prev == nil {
			// If prev == nil, we're at the beginning of the list.
			c.head = curr
		} else {
			// Otherwise, insert the node in the middle of the list.
			prev.link = curr
		}

		// The next node needs to take the new node's delta into account. Note
		// that the next node delta will alwasy be less than the current delta.
		if next != nil {
			next.delta -= delta
		}
	}

	// Add the event to the curr node.
	curr.events[t] = struct{}{}
	c.nodes[t] = curr
}

// Unschedule removes an event from the queue.
func (c *Clock[T]) Unschedule(t T) {
	if n, ok := c.nodes[t]; ok {
		delete(n.events, t)
		delete(c.nodes, t)
	}
}

// Tick advances the clock by one and pops any events with non-positive delta.
func (c *Clock[T]) Tick() []T {
	// Nothing to do if there aren't any scheduled events.
	if c.head == nil {
		return nil
	}

	// Since all deltas are relative, this decrements the entire clock.
	c.head.delta -= 1

	// Pop events from all nodes with non-positive delta.
	var events []T
	for c.head != nil && c.head.delta <= 0 {
		for t := range c.head.events {
			events = append(events, t)
		}
		c.head = c.head.link
	}

	// Cleanup pointers to discarded nodes.
	for _, t := range events {
		delete(c.nodes, t)
	}

	return events
}
