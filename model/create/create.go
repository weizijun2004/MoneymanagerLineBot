package create

import (
	"context"
	"fmt"
	"log"
	"sync"

	"moneyLineBot/model/connect"
	"moneyLineBot/model/structType"

	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// CreateGroup creates a new group with the given group name and members.
// notice: there have no ini for mongoDatabaseClient and CollectionClient
func CreateGroup(groupName string, members []string) error {
	var err error

	groupCollection, err := connect.ConnectCollectionClient(groupName)
	if err != nil {
		return err
	}

	var existEvents structType.ExistEvents
	_, err = groupCollection.InsertOne(context.TODO(), existEvents)
	if err != nil {
		log.Printf("failed to add exist events array to group")
		return fmt.Errorf("failed to add exist events array to group")
	}

	if len(members) == 0 {
		log.Printf("no members in group")
		return fmt.Errorf("no members in group")
	}
	for _, name := range members {
		var member structType.Member
		member.MemberName = name
		_, err = groupCollection.InsertOne(context.TODO(), member)
		if err != nil {
			log.Printf("failed to add \"%s\" to group", name)
			return fmt.Errorf("failed to add exist events array to group")
		}
	}
	return nil
}

// writes an event for a members in a group.
func writeEvent(group *mongo.Collection, member string, eventMembers []string, eventName string, wg *sync.WaitGroup) error {
	defer wg.Done()
	filter := bson.M{"name": member}
	var user structType.Member
	err := group.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to find user \"%s\" in group \"%s\", Error: %v", member, group.Name(), err)
	}
	var memberTemp structType.Event
	memberTemp.EventName = eventName
	for _, name := range eventMembers {
		if name != user.MemberName { // means not the user himself
			memberTemp.MemebersPay[name] = 0
		}
	}
	user.Events[eventName] = memberTemp
	group.UpdateOne(context.TODO(), filter, user)
}

// CreateEvent creates a new event for a group with the given event members and event name.
func CreateEvent(groupName string, eventMembers []string, eventName string) {
	var wg sync.WaitGroup
	group := connect.CollectionClient[groupName]
	for _, member := range eventMembers {
		wg.Add(1)
		go writeEvent(group, member, eventMembers, eventName, &wg)
	}
	wg.Wait()
	var eventExist structType.ExistEvents
	eventExist.ExistEventsArr = append(eventExist.ExistEventsArr, eventName)
}
