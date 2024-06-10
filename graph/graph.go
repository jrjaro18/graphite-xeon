package graph

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Graph struct {
	Names []string
	Map   [][]float32
}

func Create() *Graph {
	return &Graph{
		Names: make([]string, 0),
		Map:   make([][]float32, 0),
	}
}

func (g *Graph) GetNodesNames() *[]string {
	return &g.Names
}

func (g *Graph) GetGraphMap() *[][]float32 {
	return &g.Map
}

func (g *Graph) AddNode(name string) error {
	newName := strings.ToLower(strings.TrimSpace(name))
	if (newName=="") {
		return errors.New("name cannot be an empty string")
	}
	for _, v := range g.Names {
		if v == newName {
			return errors.New("node with this name already exists")
		}
	}
	g.Names = append(g.Names, newName)

	arr := make([]float32, len(g.Names))
	for i := range arr {
		arr[i] = 1
	}

	g.Map = append(g.Map, arr);
	for i := range g.Map {
		if i == len(g.Names)-1 {
			break
		}
		g.Map[i] = append(g.Map[i], 1)
	}
	return nil
}

func (g *Graph) UpdateGraphFromUserAction(userId primitive.ObjectID, postId primitive.ObjectID, action int8, timeSpent uint8) error {

	return nil
}

func (g *Graph) UpdateGraphFromPostAction(postId primitive.ObjectID, features []string) error {

	return nil
}