package create

import (
	"context"
	"log"
	"sync"

	"moneyLineBot/model/connect"
	"moneyLineBot/model/structType"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateGroup creates a new group with the given group name and members.
func CreateGroup(groupName string, member []string) {
	connect.ConnectCollectionClient(groupName)
	groupCollection := connect.CollectionClient[groupName]
	var eventExist structType.EventExist
	groupCollection.InsertOne(context.TODO(), eventExist)
	for _, name := range member {
		var member structType.User
		member.MemberName = name
		groupCollection.InsertOne(context.TODO(), member)
	}
}

// writeEvent writes an event for a specific member in a group.
func writeEvent(member string, groupName string, eventMembers []string, eventName string, wg *sync.WaitGroup) {
	defer wg.Done()
	filter := bson.M{"name": member}
	var user structType.User
	err := connect.CollectionClient[groupName].FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	var memberTemp structType.Event
	memberTemp.EventName = eventName
	for _, name := range eventMembers {
		if name != user.MemberName { // means not the user himself
			memberTemp.MemeberPay[name] = 0
		}
	}
	user.EventAttend[eventName] = memberTemp
	connect.CollectionClient[groupName].UpdateOne(context.TODO(), filter, user)
}

// CreateEvent creates a new event for a group with the given event members and event name.
func CreateEvent(groupName string, eventMembers []string, eventName string) {
	var wg sync.WaitGroup
	for _, member := range eventMembers {
		wg.Add(1)
		go writeEvent(member, groupName, eventMembers, eventName, &wg)
	}
	wg.Wait()
	var eventExist structType.EventExist
	eventExist.ExistEventArr = append(eventExist.ExistEventArr, eventName)
}
