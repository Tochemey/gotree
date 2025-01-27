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

// Node defines the structure of a Node that can be added to the Tree.
// A Node is a fundamental unit in the Tree data structure and represents
// a single item with a unique identifier and associated value. Each Node
// can be connected to other Nodes, forming a tree-like structure.
//
// The Node interface is generic, meaning it can hold values of any type
// (specified by the type parameter `T`). This makes the Tree structure
// flexible and reusable for various use cases, where the value type of the
// Nodes can vary depending on the application.
//
// Methods:
//   - ID: Returns the unique identifier for the Node, which distinguishes it
//     from other Nodes in the Tree. The ID is typically a string and must be
//     unique for each Node within the Tree.
//   - Value: Returns the value associated with the Node. This is the actual
//     data the Node holds, and its type is defined by `T`, allowing the Node
//     to hold any type of data (e.g., string, int, custom structs).
//
// Example usage:
//
//	type MyNode struct {
//	    id    string
//	    value string
//	}
//
//	func (n *MyNode) ID() string {
//	    return n.id
//	}
//
//	func (n *MyNode) Value() string {
//	    return n.value
//	}
//
//	node := &MyNode{id: "1", value: "Root Node"}
//	fmt.Println("Node ID:", node.ID())   // Output: "1"
//	fmt.Println("Node Value:", node.Value()) // Output: "Root Node"
//
// Notes:
//   - Implementations of the Node interface must ensure that the ID() method
//     returns a unique identifier for each Node. The Tree data structure relies
//     on this uniqueness to maintain correct relationships between Nodes.
//   - The Value() method can return any type of data, making the Node versatile
//     for different use cases. The type `T` can be a primitive type, a struct,
//     or even another collection type.
type Node[T any] interface {
	// ID returns the unique identifier of the Node.
	// The ID must be unique within the context of the Tree.
	ID() string
	// Value returns the value associated with the Node.
	// The type of the value is defined by the generic type T.
	Value() T
}
