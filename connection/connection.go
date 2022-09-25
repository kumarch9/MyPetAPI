package connection

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// var (
// 	ctxTime, _ = context.WithTimeout(context.TODO(), 20*time.Second)
// )

//const mongoURI = "mongodb+srv://root:12345@cluster0.wlbxjpe.mongodb.net/?retryWrites=true&w=majority"

const mongoURI = "mongodb://localhost:27017"

func ConnectionDb() *mongo.Client {
	fmt.Println("connection string is:", mongoURI)

	opt := options.Client().ApplyURI(mongoURI)
	client, errCon := mongo.Connect(context.Background(), opt)
	if errCon != nil {
		log.Println("errCon : ", errCon)
		log.Fatal(errCon)
	}

	if errPing := client.Ping(context.Background(), readpref.Primary()); errPing != nil {
		log.Println("errPing  : ", errPing)
		log.Fatal(errPing)
		//os.Exit(1)
	}

	log.Println("Database is connected to server")
	return client
}

func GetDBCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
