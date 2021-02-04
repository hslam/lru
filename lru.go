// Copyright (c) 2021 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package lru implements an LRU cache.
package lru

type node struct {
	key   interface{}
	value interface{}
	prev  *node
	next  *node
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

// New return a new LRU.
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
	for k := range l.nodes {
		delete(l.nodes, k)
	}
	l.root.next = l.root
	l.root.prev = l.root
}

// Set sets the value for a key.
func (l *LRU) Set(key, value interface{}) {
	if n, ok := l.nodes[key]; ok {
		n.value = value
		if n != l.root.next {
			l.move(n, l.root)
		}
	} else {
		if len(l.nodes)+1 > l.capacity {
			back := l.root.prev
			l.delete(back)
		}
		n := &node{key: key, value: value}
		l.nodes[key] = n
		l.insert(n, l.root)
	}
}

// Get returns the value for a key.
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	var n *node
	if n, ok = l.nodes[key]; ok {
		if n != l.root.next {
			l.move(n, l.root)
		}
		value = n.value
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
	if l.free != nil {
		l.free(n.key, n.value)
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

func (l *LRU) move(n, at *node) {
	if n != at {
		l.remove(n)
		l.insert(n, at)
	}
}
