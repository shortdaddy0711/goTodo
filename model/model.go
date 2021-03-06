package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/shortdaddy0711/goTodo/connection"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Completed bool               `json:"completed" bson:"completed"`
	SessionID string             `json:"sessionid,omitempty" bson:"sessionid,omitempty"`
	CreatedAt time.Time          `json:"createdat" bson:"createdat"`
}

func GetTodos(sessionId string) []*Todo {

	list := []*Todo{}

	collection := connection.ConnectDB()

	fmt.Println(sessionId)

	cur, err := collection.Find(context.TODO(), bson.M{"sessionid": sessionId})

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var todo *Todo

		err := cur.Decode(&todo)

		if err != nil {
			log.Fatal(err)
		}

		list = append(list, todo)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return list
}

func GetTodo(id string) *Todo {

	var todo *Todo

	collection := connection.ConnectDB()

	_id, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&todo)

	if err != nil {
		log.Fatal(err)
	}

	return todo
}

func AddTodo(r *http.Request, sessionId string) *Todo {
	var todo *Todo

	_ = json.NewDecoder(r.Body).Decode(&todo)

	todo.SessionID = sessionId
	todo.Completed = false
	todo.CreatedAt = time.Now()

	collection := connection.ConnectDB()

	result, err := collection.InsertOne(context.TODO(), todo)

	if err != nil {
		log.Fatal(err)
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)

	return todo
}

func RemoveTodo(id string) bool {
	mongoId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": mongoId}

	collection := connection.ConnectDB()

	_, err := collection.DeleteOne(context.TODO(), filter)

	return err == nil
}

func CompleteTodo(id string, complete string) bool {
	mongoId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": mongoId}

	completed := complete == "true"

	update := bson.M{"$set": bson.M{"completed": completed}}

	collection := connection.ConnectDB()

	var todo *Todo

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&todo)

	return err == nil
}
