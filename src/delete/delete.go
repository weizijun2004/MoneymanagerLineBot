package delete

import (
	"context"
	"log"

	"moneyLineBot/src/connect"
	"moneyLineBot/src/structType"

	"go.mongodb.org/mongo-driver/bson"
)

func findEvent(eventName string, groupName string, memberName string) {
	membersFilter := bson.M{"UserName": memberName}
	var user structType.User
	err := connect.CollectionClient[groupName].FindOne(context.TODO(), membersFilter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	delete(user.EventAttend, eventName)
	connect.CollectionClient[groupName].UpdateOne(context.TODO(), membersFilter, user)
}

func EventDelete(eventName string, groupName string, members []string) {
	for _, user := range members {
		go findEvent(eventName, groupName, user)
	}
}

func CollectionDelete(collectionName string) {
	delete(connect.CollectionClient, collectionName)
}
