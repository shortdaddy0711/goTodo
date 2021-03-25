package connection

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func ConnectDB() *mongo.Collection {
	uri := goDotEnvVariable("MONGODB_URI")
	db := goDotEnvVariable("DATABASE")
	col := goDotEnvVariable("COLLECTION")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to MongoDB!")
	// defer client.Disconnect(context.TODO())

	collection := client.Database(db).Collection(col)

	return collection
}
