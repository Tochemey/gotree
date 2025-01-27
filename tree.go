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
	"sort"
	"sync"
	"sync/atomic"
)

// TODO: make this configurable
const numShards = 64

// Tree defines and implements a thread-safe, flexible Tree-like data structure.
//
// This Tree allows Nodes to have an arbitrary number of children, providing a
// hierarchical structure for organizing data. Each Node in the Tree can hold
// a value of any type, and the Tree provides methods for adding, removing,
// and querying Nodes, as well as retrieving the root and all descendants.
//
// The Tree is thread-safe, meaning that it can be accessed and modified concurrently
// by multiple goroutines without causing data races. The underlying implementation
// ensures that all operations (e.g., adding/removing Nodes, querying) are safe to be
// used in multi-threaded contexts.
//
// Key Features:
//   - Each Node can have an arbitrary number of children, making the Tree suitable
//     for representing complex hierarchical relationships.
//   - The Tree supports operations to find, add, delete, and traverse Nodes.
//   - It includes thread-safety guarantees for concurrent access and modification.
//
// The Tree structure can be used to represent many types of hierarchical data,
// such as organizational charts, file systems, or tree-based decision structures.
//
// Example usage:
//
//	tree := NewTree[string]()
//	root := tree.Root()  // Get the root node
//	node := &MyNode{id: "node1", value: "Some Value"}
//	tree.Add(node, root) // Add node as a child to the root
//	descendants := tree.Descendants(root) // Retrieve descendants of root
//
// Notes:
// - The Tree structure is optimized for efficient addition and traversal of Nodes.
// - Use thread-safe methods (e.g., Add, Delete) when performing operations in concurrent environments.
type Tree[T any] struct {
	nodes      ShardedMap
	parents    ShardedMap
	nodesPool  *sync.Pool
	valuesPool *sync.Pool
	size       atomic.Int64
	// rootNode represents the tree root node
	// and there can only one root node
	rootNode *treeNode[T]
}

// Add inserts a given Node into the Tree with the specified parent Node.
//
// The `Add` method adds the `node` as a child of the specified `parent` Node in the Tree.
// If the `parent` is not set (i.e., nil), the given `node` is treated as the root Node of
// the Tree. Only one root Node can exist in the Tree; attempting to add a second root Node
// will result in an ErrInvalidOperation error.
//
// Parameters:
//   - node: The Node[T] to be added to the Tree. The Node must implement the `Node` interface,
//     providing a unique identifier via `ID()` and an associated value via `Value()`.
//   - parent: The Node[T] under which the `node` will be added as a child. If `parent` is nil,
//     the `node` will be set as the root Node of the Tree.
//
// Returns:
// - err: An error indicating the outcome of the operation. Possible values:
//   - nil: The Node was successfully added to the Tree.
//   - ErrInvalidOperation: Attempt to add a second root Node, which is not allowed.
//   - ErrParentNodeNotFound: The specified parent Node does not exist in the Tree.
//
// Notes:
//   - The `node` being added must have a unique ID not already present in the Tree.
//   - Adding a Node to a non-existent parent or a parent not currently part of the Tree will
//     result in an error.
//   - The Tree structure is updated to reflect the addition.
//
// Example usage:
//
//	tree := NewTree[string]()
//
//	root := NewNode("rootID", "rootValue")
//	err := tree.Add(root, nil) // Add root Node
//	if err != nil {
//	    log.Fatal("Error adding root:", err)
//	}
//
//	child := NewNode("childID", "childValue")
//	err = tree.Add(child, root) // Add child Node under root
//	if err != nil {
//	    log.Fatal("Error adding child:", err)
//	}
//
//	fmt.Println("Tree structure updated successfully")
func (x *Tree[T]) Add(node, parent Node[T]) (err error) {
	var (
		parentNode *treeNode[T]
		ok         bool
	)

	// check whether the node to be added is a root node
	if parent == nil && x.rootNode != nil {
		return ErrInvalidOperation
	}

	// check parent node
	if parent != nil {
		parentNode, ok = x.getNode(parent.ID())
		if !ok || parentNode == nil {
			return ErrParentNodeNotFound
		}
	}

	// get a node from the nodes pool
	childNode := x.nodesPool.Get().(*treeNode[T])
	childNode.ID = node.ID()
	val := x.valuesPool.Get().(*value[T])
	val.data = node

	// store the value atomically in the node
	childNode.SetValue(val)

	// store the node in the tree
	x.nodes.Store(node.ID(), childNode)

	// add the given node to the parent descendants
	// and update the ancestors hierarchy
	if parentNode != nil {
		parentNode.Descendants.Append(childNode)
		x.updateAncestors(parent.ID(), node.ID())
	}

	// only set the root node when parent is nil
	if parentNode == nil {
		// set the given node as root node
		x.rootNode = childNode
	}

	// increase the size
	x.size.Add(1)
	return nil
}

// Ancestors retrieves all the ancestor Nodes of a given Node in the Tree sorted by the ID.
//
// An ancestor of a Node is any Node located on the path from the root of the Tree
// to the specified Node (excluding the Node itself). The ancestors are returned
// in order, starting from the root and progressing toward the specified Node.
//
// Parameters:
// - node: The Node[T] for which the ancestors are to be retrieved.
//
// Returns:
//   - ancestors: A slice of Nodes[T] representing the ancestors of the specified Node.
//     If no ancestors exist (e.g., the Node is the root), this will be an empty slice.
//   - ok: A boolean indicating whether the operation was successful (true) or not (false).
//     It returns false if the given Node does not exist in the Tree.
//
// Notes:
//   - If the provided Node does not belong to the Tree, the method will return an empty
//     slice and `ok` will be false.
//   - This method does not modify the Tree.
//
// Example usage:
//
//	tree := NewTree[string]()
//	child, _ := tree.Find("childKey")
//
//	ancestors, ok := tree.Ancestors(child)
//	if ok {
//	    fmt.Println("Ancestors:", ancestors)
//	} else {
//	    fmt.Println("Node not found or no ancestors exist")
//	}
func (x *Tree[T]) Ancestors(node Node[T]) (ancestors []Node[T], ok bool) {
	ancestorIDs, ok := x.getAncestors(node.ID())
	if !ok {
		return nil, false
	}

	for _, ancestorID := range ancestorIDs {
		if ancestor, ok := x.getNode(ancestorID); ok {
			ancestors = append(ancestors, ancestor.GetValue())
		}
	}

	// sort the ancestors
	sort.SliceStable(ancestors, func(i, j int) bool {
		return ancestors[i].ID() < ancestors[j].ID()
	})
	return ancestors, true
}

// ParentAt retrieves the ancestor Node of the given Node at the specified level.
// The level determines the "distance" to the ancestor:
//   - A level of 0 returns the immediate parent (direct parent).
//   - A level of 1 returns the grandparent.
//   - A level of 2 returns the great-grandparent, and so on.
//
// Parameters:
//   - node: The Node[T] whose ancestor is being queried.
//   - level: A uint specifying the level of ancestry to retrieve.
//     0 represents the direct parent, 1 the grandparent, and so on.
//
// Returns:
//   - parent: The ancestor Node[T] at the given level.
//   - ok: A boolean indicating whether the ancestor was found.
//     Returns false if the level is invalid or there is no ancestor at the specified level.
//
// Usage Example:
//
//	parent, ok := tree.ParentAt(currentNode, 1)
//	if ok {
//	    fmt.Printf("Found ancestor at level 1: %+v\n", parent)
//	} else {
//	    fmt.Println("No ancestor found at the specified level")
//	}
func (x *Tree[T]) ParentAt(node Node[T], level uint) (parent Node[T], ok bool) {
	ancestor, ok := x.ancestorAt(node, int(level))
	if !ok {
		return nil, false
	}
	return ancestor.GetValue(), true
}

// Descendants retrieves all the descendant Nodes of a given Node in the Tree sorted by the ID.
//
// Descendants of a Node include all Nodes that are directly or indirectly
// connected as children of the specified Node. The descendants are returned
// in a depth-first order, starting with the immediate children and progressing
// through the subtree.
//
// Parameters:
//   - node: The Node[T] for which the descendants are to be retrieved. This Node
//     must exist in the Tree for the operation to succeed.
//
// Returns:
//   - descendants: A slice of Nodes[T] representing all the descendants of the
//     specified Node. If the Node has no descendants, this will be an empty slice.
//   - ok: A boolean indicating whether the operation was successful (true) or not
//     (false). Returns false if the specified Node does not exist in the Tree.
//
// Notes:
//   - This method does not include the specified Node itself in the results.
//   - If the provided Node does not belong to the Tree, the method will return an
//     empty slice and `ok` will be false.
//   - The Tree remains unchanged by this operation.
//
// Example usage:
//
//	tree := NewTree[string]()
//	root := tree.Root()
//	child1 := tree.Add(NewNode("child1", "child"), root)
//	child2 := tree.Add(NewNode("child1", "child"), root)
//	grandchild := tree.Add(NewNode("grandchildKey", "grandchildValue") child1)
//
//	descendants, ok := tree.Descendants(root)
//	if ok {
//	    fmt.Println("Descendants of root:", descendants)
//	} else {
//	    fmt.Println("Node not found")
//	}
func (x *Tree[T]) Descendants(node Node[T]) (descendants []Node[T], ok bool) {
	treeNode, ok := x.getNode(node.ID())
	if !ok {
		return nil, false
	}

	treeNodes := collectDescendants(treeNode)
	for _, treeNode := range treeNodes {
		descendants = append(descendants, treeNode.GetValue())
	}

	// sort the ancestors
	sort.SliceStable(descendants, func(i, j int) bool {
		return descendants[i].ID() < descendants[j].ID()
	})

	return descendants, true
}

// Delete removes the specified Node from the Tree.
//
// If the given Node exists in the Tree, it will be removed along with all its
// descendant Nodes (if any). If the Node does not exist, the method returns an
// ErrNotFound error.
//
// Parameters:
//   - node: The Node[T] to be removed from the Tree. This Node must already exist
//     in the Tree for the operation to succeed.
//
// Returns:
// - err: An error indicating the outcome of the operation. Possible values:
//   - nil: The Node was successfully deleted.
//   - ErrNotFound: The specified Node does not exist in the Tree.
//
// Notes:
//   - This operation will remove the entire subtree rooted at the specified Node.
//     Use with caution if the Node has descendants.
//   - The Tree's structure will be updated to ensure consistency after the deletion.
//   - If the Node being deleted is the root of the Tree, the Tree will be emptied.
//
// Example usage:
//
//	tree := NewTree[string]()
//	child, _ := tree.Find("childKey")
//
//	err := tree.Delete(child)
//	if err != nil {
//	    if errors.Is(err, ErrNotFound) {
//	        fmt.Println("Node not found")
//	    } else {
//	        fmt.Println("Failed to delete Node:", err)
//	    }
//	} else {
//	    fmt.Println("Node deleted successfully")
//	}
func (x *Tree[T]) Delete(node Node[T]) (err error) {
	n, ok := x.getNode(node.ID())
	if !ok {
		return ErrNotFound
	}

	// remove the node from its parent's Children slice
	if ancestors, ok := x.parents.Load(node.ID()); ok && len(ancestors.([]string)) > 0 {
		parentID := ancestors.([]string)[0]
		if parent, found := x.getNode(parentID); found {
			children := filterOutChild(parent.Descendants, node.ID())
			parent.Descendants.Reset()
			parent.Descendants.AppendMany(children.Items()...)
		}
	}

	// recursive function to delete a node and its descendants
	var deleteChildren func(n *treeNode[T])
	deleteChildren = func(n *treeNode[T]) {
		for index, child := range n.Descendants.Items() {
			n.Descendants.Delete(index)
			deleteChildren(child)
		}
		// delete node from maps and pool
		x.nodes.Delete(n.ID)
		x.parents.Delete(n.ID)
		x.nodesPool.Put(n)
		x.size.Add(-1)
	}

	deleteChildren(n)
	return nil
}

// Find searches for a Node in the Tree with the specified key.
//
// If a Node with the given key exists in the Tree, it returns the Node and
// a boolean value of true. If the Node does not exist, it returns an empty Node,
// a boolean value of false, and an ErrNotFound error.
//
// Parameters:
// - key: A string representing the unique identifier of the Node to be searched.
//
// Returns:
// - item: The Node[T] associated with the given key if found. If not found, this will be an empty Node.
// - ok: A boolean indicating whether the Node was found (true) or not (false).
//
// Note:
// - The method will not modify the Tree. It performs a read-only search operation.
// - If the Node is not found, ensure to handle the ErrNotFound error appropriately.
//
// Example usage:
//
//	tree := NewTree[string]()
//	node, ok := tree.Find("exampleKey")
//	if ok {
//	    fmt.Println("Node found:", node)
//	} else {
//	    fmt.Println("Node not found")
//	}
func (x *Tree[T]) Find(key string) (item Node[T], ok bool) {
	treeNode, ok := x.getNode(key)
	if !ok {
		return nil, false
	}
	return treeNode.GetValue(), true
}

// Root returns the root Node of the Tree.
//
// The root Node is the top-most Node in the Tree, from which all other Nodes
// (if any) descend. If the Tree is empty (i.e., no root has been set), the
// method will return a nil Node.
//
// Returns:
//   - Node[T]: The root Node of the Tree. If the Tree has no root (i.e., it is empty),
//     this will return nil.
//
// Notes:
//   - The root Node is unique and serves as the starting point for traversing the Tree.
//   - If the Tree has not yet been populated with any Nodes, calling this method
//     will return nil, which can be used to check if the Tree is empty or uninitialized.
//
// Example usage:
//
//	tree := NewTree[string]()
//	root := tree.Root()
//	if root == nil {
//	    fmt.Println("The Tree has no root")
//	} else {
//	    fmt.Println("Root node:", root)
//	}
func (x *Tree[T]) Root() Node[T] {
	return x.rootNode.GetValue()
}

// Size returns the current number of Nodes in the Tree.
//
// This method calculates and returns the total number of Nodes that have been
// added to the Tree, including the root Node and all its descendants. It provides
// a quick way to determine the size of the Tree, which can be useful for operations
// such as balancing, resource allocation, or checking whether the Tree contains any Nodes.
//
// Returns:
// - int: The current number of Nodes in the Tree. This value will be 0 if the Tree is empty.
//
// Notes:
//   - The size includes all Nodes in the Tree, regardless of their position (root, child, etc.).
//   - This method operates in constant time (O(1)) if the size is tracked, or linear time (O(n))
//     if the size is calculated by traversing all Nodes. The specific implementation
//     depends on whether the Tree structure maintains a size counter internally.
//
// Example usage:
//
//	tree := NewTree[string]()
//	fmt.Println("Tree size:", tree.Size())  // Output: 0
//
//	tree.Add(&MyNode{id: "1", value: "Root Node"}, nil)
//	tree.Add(&MyNode{id: "2", value: "Child Node"}, tree.Root())
//	fmt.Println("Tree size:", tree.Size())  // Output: 2
func (x *Tree[T]) Size() int64 {
	return x.size.Load()
}

// Reset resets the Tree to its initial state, removing all Nodes.
//
// This method clears the Tree, effectively removing the root node, all its
// descendants, and any other Nodes that were added. After calling Reset,
// the Tree will be empty, and subsequent operations (such as finding,
// adding, or removing Nodes) will operate on a newly reset Tree.
//
// Use cases for this method include situations where the Tree needs to be
// cleared and reinitialized, such as when reloading data, or when the Tree
// needs to be reset for a new operation or test case.
//
// Notes:
//   - This method does not change the underlying Tree structure itself,
//     but rather removes all Nodes, setting the Tree back to an empty state.
//   - Any references to Nodes before calling Reset will become invalid once
//     the Nodes are removed, so be cautious when retaining pointers to Nodes
//     that may be cleared.
//
// Example usage:
//
//	tree := NewTree[string]()
//	tree.Add(&MyNode{id: "1", value: "Root Node"}, nil)
//	fmt.Println("Before reset, tree size:", tree.Size()) // Output: 1
//
//	tree.Reset()
//	fmt.Println("After reset, tree size:", tree.Size()) // Output: 0
func (x *Tree[T]) Reset() {
	x.nodes.Reset()   // Reset nodes map
	x.parents.Reset() // Reset parents map
	x.size.Store(0)
}

// Nodes retrieves all the Nodes present in the Tree.
//
// This method returns a slice of all Nodes in the Tree, including the root Node.
// The Nodes are returned in an unspecified order without any hierarchy.
// If the Tree is empty, an empty slice will be returned.
//
// Returns:
//   - []Node[T]: A slice containing all the Nodes in the Tree. If the Tree is empty,
//     this will be an empty slice with no elements.
//
// Notes:
//   - The returned slice contains all Nodes in the Tree
//   - The order of Nodes in the slice is not guaranteed, and can vary based on
//     internal Tree implementation or traversal method used.
//
// Example usage:
//
//	tree := NewTree[string]()
//
//	nodes := tree.Nodes()
//	fmt.Println("All nodes in the tree:")
//	for _, node := range nodes {
//	    fmt.Println("Node ID:", node.ID(), "Value:", node.Value())
//	}
func (x *Tree[T]) Nodes() []Node[T] {
	var nodes []Node[T]
	x.nodes.Range(func(_, value any) bool {
		node := value.(*treeNode[T])
		nodes = append(nodes, node.GetValue())
		return true
	})
	return nodes
}

// NewTree creates and initializes a new instance of a Tree.
//
// This function returns a pointer to a newly created Tree, which is empty by default
// and ready to have Nodes added to it. The Tree can hold Nodes of any type, as the
// type parameter `T` is specified during initialization. The Tree structure allows
// for efficient organization and management of hierarchical data, supporting
// operations such as adding, removing, and querying Nodes.
//
// The returned Tree is empty initially, meaning there is no root Node or child Nodes.
// You can add a root Node and subsequent child Nodes using the Add method.
//
// Returns:
// - *Tree[T]: A pointer to a new, empty Tree instance.
//
// Example usage:
//
//	tree := NewTree[string]() // Create a new Tree that holds string values
//	fmt.Println("Tree size:", tree.Size()) // Output: 0 (Tree is empty)
//
//	rootNode := &MyNode{id: "root", value: "Root Node"}
//	tree.Add(rootNode, nil) // Add the root node to the Tree
//	fmt.Println("Tree size:", tree.Size()) // Output: 1 (Tree has one node)
//
// Notes:
//   - The Tree is initialized without any nodes. It must be populated with Nodes using
//     the Add method or other Tree methods.
//   - The Tree can handle nodes of any type, allowing flexible use cases for different data types.
func NewTree[T any]() *Tree[T] {
	return &Tree[T]{
		nodes:   NewShardedMap(numShards),
		parents: NewShardedMap(numShards),
		nodesPool: &sync.Pool{
			New: func() any {
				return &treeNode[T]{
					Descendants: NewSlice[*treeNode[T]](),
				}
			},
		},
		valuesPool: &sync.Pool{
			New: func() any {
				return new(value[T])
			},
		},
	}
}

func (x *Tree[T]) getNode(id string) (*treeNode[T], bool) {
	value, ok := x.nodes.Load(id)
	if !ok {
		return nil, false
	}
	node, ok := value.(*treeNode[T])
	return node, ok
}

// getAncestors returns the list of ancestor nodes
func (x *Tree[T]) getAncestors(id string) ([]string, bool) {
	if value, ok := x.parents.Load(id); ok {
		return value.([]string), true
	}
	return nil, false
}

// updateAncestors updates the parent/ancestor relationships.
func (x *Tree[T]) updateAncestors(parentID, childID string) {
	if ancestors, ok := x.getAncestors(parentID); ok {
		x.parents.Store(childID, append([]string{parentID}, ancestors...))
	} else {
		x.parents.Store(childID, []string{parentID})
	}
}

// ancestorAt retrieves the ancestor at the specified level (0 for parent, 1 for grandparent, etc.)
func (x *Tree[T]) ancestorAt(node Node[T], level int) (*treeNode[T], bool) {
	ancestors, ok := x.getAncestors(node.ID())
	if ok && len(ancestors) > level {
		return x.getNode(ancestors[level])
	}
	return nil, false
}

// collectDescendants collects all the descendants and grand children
func collectDescendants[T any](node *treeNode[T]) []*treeNode[T] {
	output := NewSlice[*treeNode[T]]()
	var recursive func(*treeNode[T])
	recursive = func(currentNode *treeNode[T]) {
		for _, child := range currentNode.Descendants.Items() {
			output.Append(child)
			recursive(child)
		}
	}
	recursive(node)
	return output.Items()
}

// filterOutChild removes the node with the given ID from the Children slice.
func filterOutChild[T any](children *Slice[*treeNode[T]], childID string) *Slice[*treeNode[T]] {
	for i, child := range children.Items() {
		if child.ID == childID {
			children.Delete(i)
			return children
		}
	}
	return children
}
