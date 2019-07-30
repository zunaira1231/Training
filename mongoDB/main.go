package main

import (

	"fmt"
	"mongoDB/config"
	"mongoDB/profile/model"
	"mongoDB/profile/repository"
	"time"

)

func main() {

	fmt.Println("Go Mongo Db")

	db, err := config.GetMongoDB()

	if err != nil {
		fmt.Println(err)
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "profile")

	//saveProfile(profileRepository)
	//updateProfile(profileRepository)
	deleteProfile(profileRepository)
	//getProfile("U2", profileRepository)
	getProfiles(profileRepository)
}

func saveProfile(profileRepository repository.ProfileRepository)  {
	var p model.Profile
	p.ID = "U3"
	p.FirstName = "Robert"
	p.LastName = "Griezmer"
	p.Email = "robert@gmail.com"
	p.Password = "123456"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	err := profileRepository.Save(&p)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Profile saved..")
	}
}

func updateProfile(profileRepository repository.ProfileRepository) {
	var p model.Profile
	p.ID = "U3"
	p.FirstName = "Wuriyanto"
	p.LastName = "Musobar"
	p.Email = "wuriyanto_musobar@gmail.com"
	p.Password = "12345678"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	err := profileRepository.Update("U3", &p)

	if err != nil {
		fmt.Println("Error detected")
		fmt.Println(err)
	} else {
		fmt.Println("Profile updated..")
	}
}

func deleteProfile(profileRepository repository.ProfileRepository) {
	err := profileRepository.Delete("U1")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Profile deleted..")
	}
}

func getProfile(id string, profileRepository repository.ProfileRepository) {
	profile, err := profileRepository.FindByID(id)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(profile.ID)
	fmt.Println(profile.FirstName)
	fmt.Println(profile.LastName)
	fmt.Println(profile.Email)
}

func getProfiles(profileRepository repository.ProfileRepository) {
	profiles, err := profileRepository.FindAll()

	if err != nil {
		fmt.Println(err)
	}

	for _, profile := range profiles{
		fmt.Println("-----------------------")
		fmt.Println(profile.ID)
		fmt.Println(profile.FirstName)
		fmt.Println(profile.LastName)
		fmt.Println(profile.Email)
	}
}
/*
// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Rest of the code will go here
}

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func main() {
	// create a new context
	ctx := context.Background()

	// create a mongo client
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:3030/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("blog")
	col := db.Collection("posts")
	// disconnects from mongo
	defer client.Disconnect(ctx)

	//////INSERT MANY/////
	// create a new context with a 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert one document
	res, err := col.InsertOne(ctx, bson.M{
		"title": "Go mongodb driver cookbook",
		"tags":  []string{"golang", "mongodb"},
		"body": `this is a long post
that goes on and on
and have many lines`,
		"comments":   1,
		"created_at": time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"new post created with id: %s",
		res.InsertedID.(primitive.ObjectID).Hex(),
	)
	// => new post created with id: 5c71caf32a346553363177ce

	//////INSERT many/////////
	// create a new context with a 10 second timeout

	res1, err := col.InsertMany(ctx, []interface{}{
		bson.M{
			"title":      "Post one",
			"tags":       []string{"golang"},
			"body":       "post one body",
			"comments":   14,
			"created_at": time.Date(2019, time.January, 10, 15, 30, 0, 0, time.UTC),
		},
		bson.M{
			"title":      "Post two",
			"tags":       []string{"nodejs"},
			"body":       "post two body",
			"comments":   2,
			"created_at": time.Now(),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted ids: %v\n", res1.InsertedIDs)
	// => inserted ids: [ObjectID("5c71ce5c6e6d43eb6e2e93be") ObjectID("5c71ce5c6e6d43eb6e2e93bf")]

//update once
	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex("5c71ce5c6e6d43eb6e2e93be")
	if err != nil {
		log.Fatal(err)
	}

	// set filters and updates
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"title": "post 2 (two)"}}

	// update document
	res2, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("modified count: %d\n", res2.ModifiedCount)
	// => modified count: 1



//update many
// set filters and updates
filter1 := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}
update1 := bson.M{"$set": bson.M{"comments": 0, "updated_at": time.Now()}}

// update documents
res3, err := col.UpdateMany(ctx, filter1, update1)
if err != nil {
log.Fatal(err)
}
fmt.Printf("modified count: %d\n", res3.ModifiedCount)
// => modified count: 17

	// set filters and updates
	filter2 := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}
	update2 := bson.M{"$set": bson.M{"comments": 0, "updated_at": time.Now()}}

	// update documents
	res4, err := col.UpdateMany(ctx, filter2, update2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("modified count: %d\n", res4.ModifiedCount)
	// => modified count: 17


//delect once
	// create ObjectID from string
	id2, err := primitive.ObjectIDFromHex("5c71ce5c6e6d43eb6e2e93be")
	if err != nil {
		log.Fatal(err)
	}

	// delete document
	res5, err := col.DeleteOne(ctx, bson.M{"_id": id2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted count: %d\n", res5.DeletedCount)
	// => deleted count: 1

//delete many
	// delete documents created older than 2 days
	filter5 := bson.M{"created_at": bson.M{
		"$lt": time.Now().Add(-2 * 24 * time.Hour),
	}}

	// update documents
	res5, err = col.DeleteMany(ctx, filter5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted count: %d\n", res5.DeletedCount)
	// => deleted count: 7
}

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Body      string             `bson:"body"`
	Tags      []string           `bson:"tags"`
	Comments  uint64             `bson:"comments"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
*/