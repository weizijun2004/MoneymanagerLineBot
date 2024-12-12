//This file used to create and close connect with database and collection

package connect

import (
	"context"
	"fmt"
	"log"
	"sync"

	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDatabaseClient *mongo.Client
var Database *mongo.Database
var CollectionClient map[string]*mongo.Collection = make(map[string]*mongo.Collection) // store group collection
var mu sync.Mutex
var Inisuccess bool = false

// to create a connect with database, and store the connect in variable
func InitializeMongoDatabaseClient() (success bool, err error) {
	clientOptions := mongoOptions.Client().ApplyURI("mongodb://localhost:27017")
	if clientOptions.Validate() != nil {
		success = false
		return success, clientOptions.Validate()
	}
	clientOptions.SetMaxPoolSize(1)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		success = false
		return success, err
	}

	mongoDatabaseClient = client
	Database = mongoDatabaseClient.Database("payManagerLineBot")
	if Database == nil {
		success = false
		return success, err
	}

	success = true
	err = nil
	return success, err
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
// func ConnectCollectionClient(groupName string) {
// 	if _, exist := CollectionClient[groupName]; !exist {
// 		CollectionClient[groupName] = Database.Collection(groupName)
// 		CollectionClient[groupName].InsertOne(context.TODO(), bson.M{"groupName": groupName})
// 	} else {
// 		CollectionClient[groupName] = Database.Collection(groupName)

// 	}
// }

func ConnectCollectionClient(groupName string) (*mongo.Collection, error) {

	if Database == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exist := CollectionClient[groupName]; !exist {
		collection := Database.Collection(groupName)
		if collection == nil {
			log.Printf("failed to get collection for group: %s", groupName)
			return nil, fmt.Errorf("failed to get collection for group: %s", groupName)
		}

		_, err := collection.InsertOne(context.TODO(), bson.M{"groupName": groupName})
		if err != nil {
			log.Printf("failed to insert initiacollectionl document for group: %s, error: %v", groupName, err)
			return nil, fmt.Errorf("failed to insert initial document for group: %s, error: %v", groupName, err)
		}

		CollectionClient[groupName] = collection
		return CollectionClient[groupName], nil
	}
	return CollectionClient[groupName], nil

}
