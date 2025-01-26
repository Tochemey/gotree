/*
 * MIT License
 *
 * Copyright (c) 2025 Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gotree

import "sync/atomic"

// value encapsulates a given treeNode value
type value[T any] struct {
	// Data represents the actual treeNode value
	data Node[T]
}

// newPidValue creates an instance of pidValue
func newValue[T any](data Node[T]) *value[T] {
	return &value[T]{data: data}
}

// Data returns the actual pidValue value
func (v *value[T]) Data() Node[T] {
	return v.data
}

// treeNode defines a treeNode on the tree
type treeNode[T any] struct {
	// ID represents the unique identifier of a treeNode
	ID string
	// Value represents the actual treeNode value
	Value atomic.Pointer[value[T]]
	// Descendants hold the list of descendants
	Descendants *safeSlice[*treeNode[T]]
}

// SetValue sets a node value
func (x *treeNode[T]) SetValue(v *value[T]) {
	x.Value.Store(v)
}

// GetValue returns the underlying value of the node
func (x *treeNode[T]) GetValue() Node[T] {
	v := x.Value.Load()
	return v.Data()
}
