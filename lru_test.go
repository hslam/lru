// Copyright (c) 2021 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package lru

import (
	"testing"
)

func TestLRU(t *testing.T) {
	capacity := 100
	free := func(key, value interface{}) {
		if key.(int) != value.(int) {
			t.Error()
		}
	}
	l := New(capacity, free)
	length := 10
	for i := 0; i < length; i++ {
		r := l.Set(i, i)
		r.Done()
		if l.root.next.key.(int) != i {
			t.Error(l.root.next.key.(int), i)
		}
		if len(l.nodes) != i+1 {
			t.Error(len(l.nodes), i+1)
		}
	}
	for i := 0; i < length; i++ {
		r := l.Set(i, i)
		r.Done()
		if l.root.next.key.(int) != i {
			t.Error(l.root.next.key.(int), i)
		}
		if len(l.nodes) != length {
			t.Error(len(l.nodes), length)
		}
	}
	for i := 0; i < length; i++ {
		_, r, _ := l.Get(i)
		r.Done()
		if l.root.next.key.(int) != i {
			t.Error(l.root.next.key.(int), i)
		}
	}
	for i := 0; i < length; i++ {
		l.Remove(i)
		if len(l.nodes) != length-i-1 {
			t.Error(len(l.nodes), length-i-1)
		}
	}
	for i := 0; i < capacity+1; i++ {
		r := l.Set(i, i)
		r.Done()
		if l.root.next.key.(int) != i {
			t.Error(l.root.next.key.(int), i)
		}
	}
	if len(l.nodes) != l.capacity {
		t.Error(len(l.nodes), l.capacity)
	}
	l.Reset()
	if len(l.nodes) != 0 {
		t.Error(len(l.nodes))
	}
}

func TestNew(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error()
		}
	}()
	New(0, nil)
}
