# gotree

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Tochemey/gotree/build.yml)
[![codecov](https://codecov.io/gh/Tochemey/gotree/graph/badge.svg?token=34NrnhK2mD)](https://codecov.io/gh/Tochemey/gotree)
[![GitHub go.mod Go version](https://badges.chse.dev/github/go-mod/go-version/Tochemey/gotree)](https://go.dev/doc/install)

Simple and thread-safe [Go](https://go.dev/) Tree library.

## Table Of Content
- [Overview](#overview)
- [Features](#features)
- [Use Cases](#use-cases)
- [Installation](#installation)
- [Methods](#methods)
    - [Note](#note) 
- [Contribution](#contribution) 

## Overview

The GoTree library is a flexible implementation of a tree-like data structure in Go. 
It allows you to organize and manage hierarchical data with ease. 
This library supports concurrent access and manipulation, making it ideal for multithreaded applications.

A tree is represented as follows:

```
root
├── node1
│   ├── subnode1
│   │   ├── sub-subnode1
│   │   └── sub-subnode2
│   ├── subnode2
│   │   └── sub-subnode3
│   └── subnode3
└── node2
    ├── subnode4
    │   ├── sub-subnode4
    │   └── sub-subnode5
    └── subnode5

```
## Features:

- **Thread-Safe**: Safe for concurrent use in multithreaded environments.
- **Flexible Node** Structure: Nodes can hold any type of data.
- **Efficient Operations**: Methods for adding, removing, and querying nodes in the tree.
- **Hierarchy Management**: Nodes can have an arbitrary number of children, enabling complex hierarchical relationships. This can come at a cost.
- **Error Handling**: Provides clear error handling when nodes cannot be found or invalid operations are attempted.

## Use Cases

- Simple decision trees implementation
- Organizational charts representation
- Hierarchical structure that requires efficient operations like adding, deleting, and finding nodes

## Installation

```bash
go get github.com/tochemey/gotree
```
## Methods

- `NewTree[T any]() *Tree[T]` - creates an instance of the Tree where T can be any golang type or user defined type.
- `Add(node, parent Node[T]) (err error)` - add a given node to the Tree. Carefully read the godoc of this method.
- `Delete(node Node[T]) (err error)` - delete a given node from the Tree and its descendants.
- `Find(key string) (item Node[T], ok bool)` - lookup a given Node on the Tree given its unique identifier.
- `Ancestors(node Node[T]) (ancestors []Node[T], ok bool)` - returns all the ancestors of a given Node.
- `ParentAt(node Node[T], level uint) (parent Node[T], ok bool)` - return the Node given parent at a given level. Carefully read the godoc of this method.
- `Descendants(node Node[T]) (descendants []Node[T], ok bool)` - return all the descendants of a given Node.
- `Root() Node[T]` - returns the root Node of the Tree.
- `Size() int64` - return the size of the Tree.
- `Reset()` - closes and resets the Tree.
- `Nodes() []Node[T]` - returns all the Nodes in the Tree.

### Note
To be able to use the `Tree` methods one need to implement the `Node[T any]` interface to define the type of Node.

## Contribution

Contributions are welcome!
The project adheres to [Semantic Versioning](https://semver.org)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

To contribute please:

- Fork the repository
- Create a feature branch
- Submit a [pull request](https://help.github.com/articles/using-pull-requests)