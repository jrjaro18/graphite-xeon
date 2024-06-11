package main

import (
	"fmt"
	"revx/db"
	g "revx/graph"
	k "revx/kafka"
	"time"
)


func main() {
	fmt.Println("Graph Based Recommendations System in Golang")
	err := db.ConnectDb()
	if err != nil {
		panic(err)
	}
	
	graph, err := g.Create("graph.txt")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Graph Created Successfully!")
	userActionsKafka := k.Initialize("test-user-actions")
	postActionsKafka := k.Initialize("test-post-actions")

	fmt.Println(graph)
	// function to close kafka and db connections
	defer func() {
		if userActionsKafka.Close() != nil {
			fmt.Println("Error closing user actions kafka")
		}
		if postActionsKafka.Close() != nil {
			fmt.Println("Error closing post actions kafka")
		}
		db.CloseDb()
	}()
	
	// Graph Upload
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fmt.Println("Uploading graph...")
			err := graph.UploadGraph("graph.txt")
			if err != nil {
				fmt.Printf("Error uploading graph: %v\n", err)
			} else {
				fmt.Println("Graph uploaded successfully")
			}
		}
	}()

	go userActionsKafka.UserActionConsumer(graph)
	postActionsKafka.PostActionConsumer(graph)


	// err := db.CreateUser("Rohan Jaiswal", []string{"f1", "race", "cars"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err := db.AddPostToUser( "6662cd284a0309874e0fe23c" ,"I just hate f1", []string{"race", "cars", "f1"}, time.Now())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// db.LikePost("6662cd284a0309874e0fe23c", "66630f44578cc53979f90020")
	//db.DisLikePost("6662d1fe3996cb18518ee8c5", "6662cf6b8bbcaa156e6fb297")
}

// user id john doe 6662cd284a0309874e0fe23c
// user id rohan jaiswal 6662d1fe3996cb18518ee8c5
// post id1 6662cf6b8bbcaa156e6fb297
// post id2 6662cfd3613d2e265140cdde