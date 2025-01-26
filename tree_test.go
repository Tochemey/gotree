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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	tree := NewTree[string]()
	root := newTestNode("root", "root")
	err := tree.Add(root, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, tree.Size())

	node1 := newTestNode("node1", "node1")
	err = tree.Add(node1, nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidOperation)

	err = tree.Add(node1, root)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, tree.Size())

	node2 := newTestNode("node2", "node2")
	err = tree.Add(node2, node1)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, tree.Size())

	ancestors, ok := tree.Ancestors(node2)
	assert.True(t, ok)
	assert.NotEmpty(t, ancestors)
	assert.Len(t, ancestors, 2)
	assert.EqualValues(t, node1.ID(), ancestors[0].ID())
	assert.EqualValues(t, root.ID(), ancestors[1].ID())

	rogue := newTestNode("rogue", "rogue")
	err = tree.Add(node1, rogue)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrParentNodeNotFound)

	ancestors, ok = tree.Ancestors(rogue)
	assert.False(t, ok)
	assert.Nil(t, ancestors)

	descendants, ok := tree.Descendants(node1)
	assert.True(t, ok)
	assert.NotEmpty(t, descendants)
	assert.Len(t, descendants, 1)
	assert.EqualValues(t, node2.ID(), descendants[0].ID())

	descendants, ok = tree.Descendants(rogue)
	assert.False(t, ok)
	assert.Empty(t, descendants)

	err = tree.Delete(rogue)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)

	actual := tree.Root()
	assert.NotNil(t, actual)
	assert.EqualValues(t, root.ID(), actual.ID())

	actual, ok = tree.Find("node1")
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.EqualValues(t, node1.ID(), actual.ID())

	actual, ok = tree.Find("node3")
	assert.False(t, ok)
	assert.Nil(t, actual)

	parent, ok := tree.ParentAt(node2, 0)
	assert.True(t, ok)
	assert.NotNil(t, parent)
	assert.EqualValues(t, node1.ID(), parent.ID())

	parent, ok = tree.ParentAt(node2, 1)
	assert.True(t, ok)
	assert.NotNil(t, parent)
	assert.EqualValues(t, root.ID(), parent.ID())

	parent, ok = tree.ParentAt(node2, 2)
	assert.False(t, ok)
	assert.Nil(t, parent)

	node3 := newTestNode("node3", "node3")
	err = tree.Add(node3, node2)
	assert.NoError(t, err)

	node4 := newTestNode("node4", "node4")
	err = tree.Add(node4, node3)
	assert.NoError(t, err)

	descendants, ok = tree.Descendants(node1)
	assert.True(t, ok)
	assert.NotEmpty(t, descendants)
	assert.EqualValues(t, 3, len(descendants))

	err = tree.Delete(node2)
	assert.NoError(t, err)

	all := tree.Nodes()
	assert.NotEmpty(t, all)
	assert.Len(t, all, 2)

	tree.Reset()
}
