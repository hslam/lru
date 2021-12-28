// Copyright (c) 2021 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package lru implements an LRU cache.
package lru

// Reference represents the reference counter.
type Reference interface {
	Done()
}

type node struct {
	key     interface{}
	value   interface{}
	cost    int
	prev    *node
	next    *node
	free    Free
	counter int64
}

type reference struct {
	n    *node
	step int64
}

func (n *node) newReference() Reference {
	return &reference{n, 1}
}

// Done decrements the reference counter by one.
func (r *reference) Done() {
	r.n.counter -= r.step
	r.step = 0
	if r.n.free != nil && r.n.counter < 0 {
		r.n.free(r.n.key, r.n.value)
	}
}

// Free is a free callback.
type Free func(key, value interface{})

// LRU represents an LRU cache.
type LRU struct {
	nodes    map[interface{}]*node
	root     *node
	capacity int
	cost     int
	free     Free
}

// New returns a new LRU cache. LRU is not thread safe.
func New(capacity int, free Free) *LRU {
	if capacity <= 0 {
		panic("non-positive capacity")
	}
	l := &LRU{
		nodes:    make(map[interface{}]*node),
		root:     &node{},
		capacity: capacity,
		free:     free,
	}
	l.Reset()
	return l
}

// Reset resets the LRU cache.
func (l *LRU) Reset() {
	for _, n := range l.nodes {
		l.delete(n)
	}
	l.root.next = l.root
	l.root.prev = l.root
}

// Resize sets the LRU cache capacity.
func (l *LRU) Resize(capacity int) {
	if capacity <= 0 {
		panic("non-positive capacity")
	}
	l.capacity = capacity
	for l.cost > l.capacity {
		back := l.root.prev
		l.delete(back)
	}
}

// Size returns the LRU cache size
func (l *LRU) Size() int {
	return l.cost
}

// Set sets the value and increments the reference counter by one for a key.
func (l *LRU) Set(key, value interface{}, cost int) (reference Reference, ok bool) {
	var n = &node{key: key, value: value, cost: cost, free: l.free}
	if cost > l.capacity {
		return n.newReference(), false
	}
	if old, ok := l.nodes[key]; ok {
		var removed bool
		var oldCost = old.cost
		for l.cost-oldCost+cost > l.capacity {
			back := l.root.prev
			if old == back {
				removed = true
				oldCost = 0
			}
			l.delete(back)
		}
		l.nodes[key] = n
		if removed {
			l.insert(n, l.root)
		} else {
			l.replace(old, n)
			old.newReference().Done()
			if n != l.root.next {
				l.move(n, l.root)
			}
		}
		l.cost = l.cost - oldCost + cost
	} else {
		for l.cost+cost > l.capacity {
			back := l.root.prev
			l.delete(back)
		}
		l.nodes[key] = n
		l.insert(n, l.root)
		l.cost += cost
	}
	n.counter++
	return n.newReference(), true
}

// Remove removes the value for a key.
func (l *LRU) Remove(key interface{}) (ok bool) {
	var n *node
	if n, ok = l.nodes[key]; ok {
		l.delete(n)
	}
	return
}

// Get returns the value and increments the reference counter by one for a key.
func (l *LRU) Get(key interface{}) (value interface{}, reference Reference, ok bool) {
	var n *node
	if n, ok = l.nodes[key]; ok {
		if n != l.root.next {
			l.move(n, l.root)
		}
		value = n.value
		n.counter++
		reference = n.newReference()
	}
	return
}

func (l *LRU) delete(n *node) {
	delete(l.nodes, n.key)
	l.remove(n)
	n.newReference().Done()
	l.cost -= n.cost
}

func (l *LRU) move(n, at *node) {
	if n != at {
		l.remove(n)
		l.insert(n, at)
	}
}

func (l *LRU) insert(n, at *node) {
	n.prev = at
	n.next = at.next
	n.prev.next = n
	n.next.prev = n
}

func (l *LRU) remove(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil
	n.prev = nil
}

func (l *LRU) replace(old, new *node) {
	new.next = old.next
	new.next.prev = new
	new.prev = old.prev
	new.prev.next = new
}
