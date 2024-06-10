package main

import (
	"fmt"
	"revx/db"
	"revx/graph"
	"revx/kafka"
	// "time"
)


func main() {
	fmt.Println("Graph Based Recommendations System in Golang")
	err := db.ConnectDb()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDb Successfully!")
	graph := graph.Create()
	fmt.Println("Graph Created Successfully!")
	userActionsKafka := kafka.Initialize("test-user-actions")
	postActionsKafka := kafka.Initialize("test-post-actions")

	defer func() {
		if userActionsKafka.Close() != nil {
			fmt.Println("Error closing user actions kafka")
		}
		if postActionsKafka.Close() != nil {
			fmt.Println("Error closing post actions kafka")
		}
		db.CloseDb()
	}()

	go kafka.UserActionConsumer(userActionsKafka, graph)
	kafka.PostActionConsumer(postActionsKafka, graph)


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