package delete

import (
	"context"
	"log"
	"sync"

	"moneyLineBot/model/connect"
	"moneyLineBot/model/structType"

	"go.mongodb.org/mongo-driver/bson"
)

func findAndDeleteEvent(eventName string, groupName string, memberName string, wg *sync.WaitGroup) {
	defer wg.Done()
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
	var wg sync.WaitGroup
	for _, user := range members {
		wg.Add(1)
		go findAndDeleteEvent(eventName, groupName, user, &wg)
	}
	wg.Wait()
}

func CollectionDelete(collectionName string) {
	delete(connect.CollectionClient, collectionName)
}
