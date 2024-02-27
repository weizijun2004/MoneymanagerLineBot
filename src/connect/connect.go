//This file used to create and close connect with database and collection

package connect

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDatabaseClient *mongo.Client
var Database *mongo.Database
var CollectionClient map[string]*mongo.Collection // store group collection

// to create a connect with database, and store the connect in variable
func InitializeMongoDatabaseClient() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetMaxPoolSize(10)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	mongoDatabaseClient = client
	Database = mongoDatabaseClient.Database("payManagerLineBot")
	return nil
}

// to close the database connect
func CloseMongoDatabaseClient() {
	if mongoDatabaseClient != nil && len(CollectionClient) == 0 {
		err := mongoDatabaseClient.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// to add a group collection in group collection map
func ConnectCollectionClient(groupName string) {
	if _, exist := CollectionClient[groupName]; exist {
		CollectionClient[groupName] = Database.Collection(groupName)
	} else {
		CollectionClient[groupName] = Database.Collection(groupName)
		CollectionClient[groupName].InsertOne(context.TODO(), bson.M{"groupName": groupName})
	}
}
