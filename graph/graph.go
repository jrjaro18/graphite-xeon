package graph

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Graph struct {
	Names []string `json:"names"`
	Map   [][]float32 `json:"map"`
}

func Create(filename string) (*Graph, error) {
	stat, err := os.Stat(filename)
	if (err != nil) {
		os.Create(filename)
		return &Graph{
			Names: make([]string, 0), 
			Map: make([][]float32, 0),
		}, nil
	}
	file, err := os.Open(filename)
	defer func() {
		file.Close()
	}()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return &Graph{
			Names: make([]string, 0),
			Map: make([][]float32, 0),
		}, nil
	}
	byteGraph, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	graph := &Graph{}
	err = json.Unmarshal(byteGraph, graph)
	if err != nil {
		return nil, err
	}
	return graph, nil
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

func (g *Graph) UploadGraph(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer func() {
		file.Close()
	}()
	if err != nil {
		return err
	}
	byteGraph, err := json.Marshal(g)
	if err != nil {	
		return err
	}
	_, err = file.Write(byteGraph)
	if err != nil {
		return err
	}
	return nil
}

func (g *Graph) UpdateGraphFromUserAction(userId primitive.ObjectID, postId primitive.ObjectID, action int8, timeSpent uint8) error {

	return nil
}

func (g *Graph) UpdateGraphFromPostAction(postId primitive.ObjectID, features []string) error {
	time.Sleep(30 * time.Second)
	return nil
}