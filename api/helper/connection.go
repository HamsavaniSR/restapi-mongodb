package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//ConnectDB Mongo Connection
func ConnectDB() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

	if err != nil {
		fmt.Println("Successfully connected and pinged.panic 1 ", err)
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println("Successfully connected and pinged.panic 2 ", err)
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("Successfully connected and pinged.panic 3 ", err)
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer cancel()
	//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//clientOptions := options.Client().SetDirect(true).ApplyURI("mongodb://localhost:27017/?compressors=disabled&gssapiServiceName=mongodb")
	/*
		TODO returns a non-nil, empty Context.
		Code should use context.TODO when it's unclear which Context to use or
		it is not yet available
		(because the surrounding function has not yet been extended
		 to accept a Context parameter).
	*/
	//	client, err := mongo.Connect(context.TODO(), clientOptions)
	//	err1 := client.Ping(context.TODO(), nil)
	/*ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err1 := client.Ping(ctx, readpref.Primary())
	fmt.Println("Ping===", err1)*/
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB client===", client)
	fmt.Println("Connected to MongoDB err===", err)
	collection := client.Database("go_rest_api").Collection("books")
	fmt.Println("Connected to MongoDB collection===", *collection)
	return collection
}

//ErrorResponse for error model
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

//GetError : To prepare Error Model
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	// Marshal returns the JSON encoding of v.
	message, _ := json.Marshal(response) // returns byte,error
	w.WriteHeader(response.StatusCode)
	w.Write(message)
	fmt.Println(response)
}
