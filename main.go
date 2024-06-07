package main

import (
	"fmt"
	"revx/db"
	// "time"
)


func main() {
	fmt.Println("Graph Based Recommendations System in Golang")
	db.ConnectDb()
	defer db.CloseDb()
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