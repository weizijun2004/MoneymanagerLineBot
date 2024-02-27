package create

import (
	"context"
	"log"

	"moneyLineBot/src/connect"
	"moneyLineBot/src/structType"

	"go.mongodb.org/mongo-driver/bson"
)

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

func writeEvent(member string, groupName string, eventMembers []string, eventName string) {
	filter := bson.M{"name": member}
	var user structType.User
	err := connect.CollectionClient[groupName].FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	var memberTemp structType.Event
	memberTemp.EventName = eventName
	for _, name := range eventMembers {
		if name != user.MemberName {
			memberTemp.MemeberPay[name] = 0
			user.EventAttend[eventName] = memberTemp
		}
	}
	connect.CollectionClient[groupName].UpdateOne(context.TODO(), filter, user)
}

func CreateEvent(groupName string, eventMembers []string, eventName string) {
	for _, member := range eventMembers {
		go writeEvent(member, groupName, eventMembers, eventName)
	}
}
