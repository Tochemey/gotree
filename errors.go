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

import "errors"

var (
	// ErrParentNodeNotFound is returned when attempting to add a Node to the Tree,
	// but the specified parent Node cannot be found in the Tree.
	//
	// This error indicates that the parent Node does not exist in the Tree,
	// meaning the operation cannot proceed. It is typically used in methods like
	// Add(), where the parent Node must exist for the Node to be added as a child.
	//
	// Example usage:
	//   err := tree.Add(child, nonExistentParent)
	//   if errors.Is(err, ErrParentNodeNotFound) {
	//       fmt.Println("The specified parent node does not exist.")
	//   }
	ErrParentNodeNotFound = errors.New("parent node not found")

	// ErrNotFound is returned when a specified Node cannot be found in the Tree.
	//
	// This error occurs when attempting to locate a Node (e.g., via a Find() or
	// Descendants() method) and the Node does not exist in the Tree. The error
	// serves as a signal that the Node is not part of the Tree, allowing the
	// caller to handle the case where the search fails.
	//
	// Example usage:
	//   node, ok := tree.Find("nodeID")
	//   if !ok {
	//       fmt.Println("Node not found:", ErrNotFound)
	//   }
	ErrNotFound = errors.New("node not found")

	// ErrInvalidOperation is returned when an invalid operation is attempted,
	// such as trying to add a second root Node to the Tree, which is not allowed.
	//
	// This error indicates that the Tree already has a root Node, and attempting
	// to add another root Node would violate the Tree's constraints. The error
	// helps prevent invalid state transitions in the Tree structure.
	//
	// Example usage:
	//   err := tree.Add(newRoot, nil)
	//   if errors.Is(err, ErrInvalidOperation) {
	//       fmt.Println("Cannot add a second root node:", ErrInvalidOperation)
	//   }
	ErrInvalidOperation = errors.New("invalid operation")
)
