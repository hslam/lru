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
	prev    *node
	next    *node
	free    Free
	counter int64
}

// Done decrements the reference counter by one.
func (n *node) Done() {
	n.counter--
	if n.free != nil && n.counter < 0 {
		n.free(n.key, n.value)
	}
}

// Free is a free callback.
type Free func(key, value interface{})

// LRU represents an LRU cache.
type LRU struct {
	nodes    map[interface{}]*node
	root     *node
	capacity int
	free     Free
}

// New returns a new LRU cache.
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

// Set sets the value and increments the reference counter by one for a key.
func (l *LRU) Set(key, value interface{}) Reference {
	var n *node
	var ok bool
	if n, ok = l.nodes[key]; ok {
		n.value = value
		if n != l.root.next {
			l.move(n, l.root)
		}
	} else {
		if len(l.nodes)+1 > l.capacity {
			back := l.root.prev
			l.delete(back)
		}
		n = &node{key: key, value: value, free: l.free}
		l.nodes[key] = n
		l.insert(n, l.root)
	}
	n.counter++
	return n
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
		reference = n
	}
	return
}

// Remove removes the value for a key.
func (l *LRU) Remove(key interface{}) (ok bool) {
	var n *node
	if n, ok = l.nodes[key]; ok {
		l.delete(n)
	}
	return
}

func (l *LRU) delete(n *node) {
	delete(l.nodes, n.key)
	l.remove(n)
	n.Done()
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

func (l *LRU) move(n, at *node) {
	if n != at {
		l.remove(n)
		l.insert(n, at)
	}
}
