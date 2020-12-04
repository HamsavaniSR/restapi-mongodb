package main

import (
	"context"
	"encoding/json"

	"github.com/restapi-mongodb/api/helper"
	"github.com/restapi-mongodb/api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var collection = helper.ConnectDB()

func main() {
	fmt.Println("Starting the application ")
	r := mux.NewRouter()
	r.HandleFunc("/test", test).Methods("GET")
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Hi")
	json.NewEncoder(w).Encode("This is a test")
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book) // get address of book
	fmt.Println("Create book is ", book)
	result, err := collection.InsertOne(context.TODO(), book)
	fmt.Println("Create book result==", err)
	if err != nil {
		helper.GetError(err, w)

		return
	}
	json.NewEncoder(w).Encode(result)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []models.Book
	//bson.M{} unordered representation of docs
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	// Close the cursor once finished
	/*A defer (otthivai) statement defers the execution of a function until
	the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var book models.Book
		err := cursor.Decode(&book) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(books) // encode similar to deserialize process.
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book

	var params = mux.Vars(r)
	//string to objectID conversion
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book

	var params = mux.Vars(r)
	//string to objectID conversion
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	json.NewDecoder(r.Body).Decode(&book) // get address of book
	// prepare update model
	update := bson.D{
		{"$set", bson.D{
			{"isbn", book.Isbn},
			{"title", book.Title},
			{"author", bson.D{
				{"firstname", book.Author.FirstName},
				{"lastname", book.Author.LastName},
			}},
		}},
	}
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	book.ID = id // update id
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//	var book models.Book

	var params = mux.Vars(r)
	//string to objectID conversion
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}
