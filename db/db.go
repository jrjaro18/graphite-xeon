package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collUser *mongo.Collection
var collPost *mongo.Collection
var cli *mongo.Client

// User represents a user in the system
type User struct {
	Name     string               `bson:"name"`
	Posts    []primitive.ObjectID `bson:"posts"`
	Features []string             `bson:"features"`
	AdditionalFeatures []string   `bson:"additional_features"`
}

// Post represents a post made by a user
type Post struct {
	Text        string               `bson:"text"`
	Features    []string             `bson:"features"`
	CreatedAt   time.Time            `bson:"created_at"`
	Likes       []primitive.ObjectID `bson:"likes"`
	DisLikes    []primitive.ObjectID `bson:"dislikes"`
	LikesToday  int                  `bson:"likes_today"`
	CurrentDate time.Time            `bson:"current_date"`
}

// connect to db
func ConnectDb() {
	// Connect to the database
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Get a handle for the users collection
	cli = client
	collUser = client.Database("graphite").Collection("users")
	collPost = client.Database("graphite").Collection("posts")
}

// close db connection
func CloseDb() {
	err := cli.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// create a new user
func CreateUser(username string, features []string) error {
	//check if user already exists
	var existingUser User
	err := collUser.FindOne(context.Background(), bson.M{"name": username}).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("user with name %s already exists", username)
	}

	user := User{
		Name:     username,
		Posts:    []primitive.ObjectID{},
		Features: features,
	}

	_, err = collUser.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// add a post to a user
func AddPostToUser(userId string, postText string, features []string, t time.Time) error {
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}
	//check if user exists
	var existingUser User
	err = collUser.FindOne(context.Background(), bson.M{"_id": userObjectId}).Decode(&existingUser)
	if err != nil {
		return fmt.Errorf("user with id %s does not exist", userId)
	}
	post := Post{
		Text:      postText,
		Features:  features,
		CreatedAt: t,
		Likes:     []primitive.ObjectID{},
		DisLikes:  []primitive.ObjectID{},
		CurrentDate: time.Now().Truncate(24 * time.Hour),
	}
	res, err := collPost.InsertOne(context.Background(), post)
	if err != nil {
		return err
	}
	_, err = collUser.UpdateOne(context.Background(), bson.M{"_id": userObjectId}, bson.M{"$push": bson.M{"posts": res.InsertedID}})
	if err != nil {
		return err
	}

	return nil
}

// like a post
func LikePost(userId string, postId string) error {
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}
	postObjectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return fmt.Errorf("invalid post ID format: %v", err)
	}
	//check if user exists
	var existingUser User
	err = collUser.FindOne(context.Background(), bson.M{"_id": userObjectId}).Decode(&existingUser)
	if err != nil {
		return fmt.Errorf("user with id %s does not exist", userId)
	}
	//check if post exists
	var existingPost Post
	err = collPost.FindOne(context.Background(), bson.M{"_id": postObjectId}).Decode(&existingPost)
	if err != nil {
		return fmt.Errorf("post with id %s does not exist", postId)
	}
	//  check if user has already liked the post or disliked the post
	for _, v := range existingPost.Likes {
		if v == userObjectId {
			// unliking the post
			_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$pull": bson.M{"likes": userObjectId}})
			if err != nil {
				return err
			}
			return nil
		}
	}
	for _, v := range existingPost.DisLikes {
		if v == userObjectId {
			// undislike the post
			_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$pull": bson.M{"dislikes": userObjectId}, "$push": bson.M{"likes": userObjectId}})
			if err != nil {
				return err
			}
			return updateLikesToday(postObjectId)
		}
	}

	_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$push": bson.M{"likes": userObjectId}})
	if err != nil {
		return err
	}
	return updateLikesToday(postObjectId)
}

// dislike a post
func DisLikePost(userId string, postId string) error {
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}
	postObjectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return fmt.Errorf("invalid post ID format: %v", err)
	}
	//check if user exists
	var existingUser User
	err = collUser.FindOne(context.Background(), bson.M{"_id": userObjectId}).Decode(&existingUser)
	if err != nil {
		return fmt.Errorf("user with id %s does not exist", userId)
	}
	//check if post exists
	var existingPost Post
	err = collPost.FindOne(context.Background(), bson.M{"_id": postObjectId}).Decode(&existingPost)
	if err != nil {
		return fmt.Errorf("post with id %s does not exist", postId)
	}
	//  check if user has already liked the post or disliked the post
	for _, v := range existingPost.DisLikes {
		if v == userObjectId {
			// undisliking the post
			_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$pull": bson.M{"dislikes": userObjectId}})
			if err != nil {
				return err
			}
			return nil
		}
	}
	for _, v := range existingPost.Likes {
		if v == userObjectId {
			// unlike the post
			_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$pull": bson.M{"likes": userObjectId}, "$push": bson.M{"dislikes": userObjectId}})
			if err != nil {
				return err
			}
			return nil
		}
	}

	_, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postObjectId}, bson.M{"$push": bson.M{"dislikes": userObjectId}})
	if err != nil {
		return err
	}
	return nil
}

// used in Like Post helps in handling the likes_today field 
/* known issue: if the like button is clicked multiple times by the same user on 
the same post on the same day, the likes_today field will be incremented multiple times */
// this can be fixed by adding a check in the frontend
func updateLikesToday(postId primitive.ObjectID) error {
    // Retrieve the existing post
    var existingPost Post
    err := collPost.FindOne(context.Background(), bson.M{"_id": postId}).Decode(&existingPost)
    if err != nil {
        return fmt.Errorf("post with id %s does not exist", postId.Hex())
    }

    // Check the current date
    currentDate := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
    var update bson.M
    if existingPost.CurrentDate.Truncate(24 * time.Hour).Equal(currentDate) {
        // Same day, increment likes_today
        update = bson.M{"$inc": bson.M{"likes_today": 1}}
    } else {
        // New day, reset likes_today to 1 and update current_date
        update = bson.M{"$set": bson.M{"likes_today": 1, "current_date": currentDate}}
    }

    // Apply the update
    _, err = collPost.UpdateOne(context.Background(), bson.M{"_id": postId}, update)
    if err != nil {
        return err
    }

    return nil
}