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
	var wg sync.WaitGroup
	errorChan := make(chan error, len(members))

	if len(members) == 0 {
		log.Printf("no members in group")
		return fmt.Errorf("no members in group")
	}

	groupCollection, err := connect.ConnectCollectionClient(groupName)
	if err != nil {
		return err
	}

	var newGroup structType.Group
	newGroup.GroupName = groupName
	for _, member := range members {
		newGroup.GroupMembers = append(newGroup.GroupMembers, member)
	}
	_, err = groupCollection.InsertOne(context.TODO(), newGroup)
	if err != nil {
		return err
	}

	for _, memberName := range members {
		wg.Add(1)
		go func(memberName string) {
			defer wg.Done()
			err := CreateMember(memberName, groupCollection)
			if err != nil {
				errorChan <- err
			}
		}(memberName)
	}
	wg.Wait()
	close(errorChan)
	var errors []error
	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", errors)
	}
	return nil
}

func CreateMember(memberName string, groupColl *mongo.Collection) error {
	// defer wg.Done()
	var member structType.Member
	member.MemberName = memberName
	_, err := groupColl.InsertOne(context.TODO(), member)
	if err != nil {
		return err
	}
	return nil
}

// CreateEvent creates a new event for a group with the given event members and event name.
func CreateEvent(groupName string, eventMembers []string, eventName string) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(eventMembers))
	if len(eventMembers) <= 0 {
		log.Printf("no members in group")
		return fmt.Errorf("no members in group")
	}
	groupCollection, err := connect.ConnectCollectionClient(groupName)
	if err != nil {
		return err
	}
	for _, member := range eventMembers {
		wg.Add(1)
		go func(memberName string) {
			defer wg.Done()
			localErr := writeEvent(groupCollection, memberName, eventMembers, eventName)
			if localErr != nil {
				errorChan <- localErr
			}
		}(member)
	}
	wg.Wait()
	close(errorChan)
	var errors []error
	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", errors)
	}
	return nil
}

// writes an event for a members in a group.
func writeEvent(groupColl *mongo.Collection, memberName string, eventMembers []string, eventName string) error {
	filter := bson.M{"name": memberName}
	var member structType.Member
	err := groupColl.FindOne(context.TODO(), filter).Decode(&member)
	if err != nil {
		return fmt.Errorf("failed to find user \"%s\" in group \"%s\", Error: %v", member, groupColl.Name(), err)
	}
	// var memberTemp structType.Event
	memberEvents := member.Events

	eventTemp := memberEvents[eventName]
	eventTemp.EventName = eventName
	if eventTemp.MemebersPay == nil {
		eventTemp.MemebersPay = make(map[string]int)
	}
	for _, name := range eventMembers {
		if name != member.MemberName { // means not the user himself
			eventTemp.MemebersPay[name] = 0
		}
	}
	memberEvents[eventName] = eventTemp
	member.Events = memberEvents
	_, err = groupColl.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"Events": member.Events}})
	if err != nil {
		return err
	}
	return nil
}
