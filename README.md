# gotree

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Tochemey/gotree/build.yml)
![Codecov](https://img.shields.io/codecov/c/github/tochemey/gotree)
[![GitHub go.mod Go version](https://badges.chse.dev/github/go-mod/go-version/Tochemey/gotree)](https://go.dev/doc/install)

Simple and thread-safe [Go](https://go.dev/) Tree library.

## Overview

The GoTree library is a flexible and thread-safe implementation of a tree-like data structure in Go. 
It allows you to organize and manage hierarchical data with ease. 
Nodes in the tree can have an arbitrary number of children, and the library provides various methods for querying, manipulating, and traversing the tree structure. 

This library supports concurrent access and manipulation, making it ideal for multi-threaded applications.
With this library, you can build structures such as organizational charts, decision trees, file systems, or any other hierarchical structure that requires efficient operations like adding, deleting, and finding nodes.

```
root
├── node1
│   ├── subnode1
│   │   ├── sub-subnode1
│   │   └── sub-subnode2
│   ├── subnode2
│   │   └── sub-subnode3
│   └── subnode3
├── node2
│   ├── subnode4
│   └── subnode5
├── node3
│   ├── subnode6
│   │   ├── sub-subnode4
│   │   └── sub-subnode5
│   ├── subnode7
│   └── subnode8
└── node4
    ├── subnode9
    │   └── sub-subnode6
    └── subnode10
```

## Features:

- **Thread-Safe**: Safe for concurrent use in multi-threaded environments.
- **Flexible Node Structure**: Nodes can hold any type of data.
- **Efficient Operations**: Methods for adding, removing, and querying nodes in the tree.
- **Hierarchy Management**: Nodes can have an arbitrary number of children, enabling complex hierarchical relationships.
- **Error Handling**: Provides clear error handling when nodes cannot be found or invalid operations are attempted.

## Installation

```bash
go get github.com/tochemey/gotree
```

## Contribution

Contributions are welcome!
The project adheres to [Semantic Versioning](https://semver.org)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

To contribute please:

- Fork the repository
- Create a feature branch
- Submit a [pull request](https://help.github.com/articles/using-pull-requests)