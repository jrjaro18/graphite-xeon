// graph_test.go
package graph

import (
    "testing"
)

func TestCreate(t *testing.T) {
    g := Create()
    if g == nil {
        t.Fatalf("Expected non-nil Graph instance")
    }
    if len(*g.GetNodesNames()) != 0 {
        t.Fatalf("Expected no nodes in new graph, got %d", len(*g.GetNodesNames()))
    }
    if len(*g.GetGraphMap()) != 0 {
        t.Fatalf("Expected empty graph map in new graph, got %d", len(*g.GetGraphMap()))
    }
}

func TestAddNode(t *testing.T) {
    g := Create()

    // Test adding a single node
    err := g.AddNode("Node1")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(*g.GetNodesNames()) != 1 {
        t.Fatalf("Expected 1 node, got %d", len(*g.GetNodesNames()))
    }
    if (*g.GetNodesNames())[0] != "node1" {
        t.Fatalf("Expected node name to be 'node1', got %s", (*g.GetNodesNames())[0])
    }

    // Test adding a duplicate node
    err = g.AddNode("Node1")
    if err == nil {
        t.Fatalf("Expected error for duplicate node, got nil")
    }

    // Test adding another node
    err = g.AddNode("Node2")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(*g.GetNodesNames()) != 2 {
        t.Fatalf("Expected 2 nodes, got %d", len(*g.GetNodesNames()))
    }

    // Test graph map size
    graphMap := *g.GetGraphMap()
    if len(graphMap) != 2 || len(graphMap[0]) != 2 || len(graphMap[1]) != 2 {
        t.Fatalf("Graph map size incorrect, got %v", graphMap)
    }
    if graphMap[0][1] != 1 || graphMap[1][0] != 1 {
        t.Fatalf("Graph map values incorrect, got %v", graphMap)
    }
}

func TestAddNode_EdgeCases(t *testing.T) {
    g := Create()

    // Test adding a node with leading and trailing spaces
    err := g.AddNode("  Node3  ")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if (*g.GetNodesNames())[0] != "node3" {
        t.Fatalf("Expected node name to be 'node3', got %s", (*g.GetNodesNames())[0])
    }

    // Test adding a node with an empty string
    err = g.AddNode("")
    if err == nil {
        t.Fatalf("Expected error for empty string, got nil")
    }

    // Test adding a node with only spaces
    err = g.AddNode("   ")
    if err == nil {
        t.Fatalf("Expected error for spaces only, got nil")
    }

    // Test adding a node with mixed case
    err = g.AddNode("NoDe4")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if (*g.GetNodesNames())[1] != "node4" {
        t.Fatalf("Expected node name to be 'node4', got %s", (*g.GetNodesNames())[1])
    }
}

func TestGraphMapInitialization(t *testing.T) {
    g := Create()

    g.AddNode("A")
    g.AddNode("B")
    g.AddNode("C")

    graphMap := *g.GetGraphMap()
    if len(graphMap) != 3 {
        t.Fatalf("Expected graph map size to be 3, got %d", len(graphMap))
    }

    for i := range graphMap {
        if len(graphMap[i]) != 3 {
            t.Fatalf("Expected row %d size to be 3, got %d", i, len(graphMap[i]))
        }
        for j := range graphMap[i] {
            if graphMap[i][j] != 1 {
                t.Fatalf("Expected graphMap[%d][%d] to be 1, got %f", i, j, graphMap[i][j])
            }
        }
    }
}
