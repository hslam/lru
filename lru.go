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

// LRU represents an LRU cache.
type LRU struct {
	nodes    map[interface{}]*node
	root     *node
	capacity int
}

// New return a new LRU.
func New(capacity int) *LRU {
	if capacity <= 0 {
		panic("non-positive capacity")
	}
	l := &LRU{
		nodes:    make(map[interface{}]*node),
		root:     &node{},
		capacity: capacity,
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

// Set sets a value with a key.
func (l *LRU) Set(key interface{}, value interface{}) {
	if n, ok := l.nodes[key]; ok {
		if n != l.root.next {
			l.move(n, l.root)
		}
	} else {
		n := &node{key: key, value: value}
		l.nodes[key] = n
		l.insert(n, l.root)
	}
	if len(l.nodes) > l.capacity {
		back := l.root.prev
		l.remove(back)
	}
}

// Get gets a value with a key.
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

// Remove removes a value with a key.
func (l *LRU) Remove(key interface{}) (ok bool) {
	var n *node
	if n, ok = l.nodes[key]; ok {
		delete(l.nodes, key)
		l.remove(n)
	}
	return
}

func (l *LRU) insert(n, at *node) *node {
	n.prev = at
	n.next = at.next
	n.prev.next = n
	n.next.prev = n
	return n
}

func (l *LRU) remove(n *node) *node {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil
	n.prev = nil
	return n
}

func (l *LRU) move(n, at *node) *node {
	if n != at {
		l.remove(n)
		l.insert(n, at)
	}
	return n
}
