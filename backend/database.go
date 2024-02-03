package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type UserID = uint64
type GroupID = uint64
type UnixMillis = uint64

// A service capable of persisting data.
type Database interface {
	// Creates a new user in the database.
	//
	// Returns an error if a user with the same `UserID`
	// already exists, or if the operation may have failed.
	CreateUser(User) error
	// Reads a user from the database.
	//
	// Returns a nil `*User` if no such user exists. Returns
	// an error if the operation could not be completed.
	ReadUser(UserID) (*User, error)
	// Deletes a user from the database, if it exists.
	//
	// Returns an error if the operation could not be completed.
	DeleteUser(UserID) error
	// Creates a new group in the database.
	//
	// Returns an error if a group with the same `GroupID`
	// already exists, or if the operation may have failed.
	CreateGroup(Group) error
	// Reads a group from the database.
	//
	// Returns a nil `*Group` if no such group exists. Returns
	// an error if the operation could not be completed.
	ReadGroup(GroupID) (*Group, error)
	// Updates the name of the group.
	//
	// Returns an error if the operation could not be completed.
	UpdateGroupName(GroupID, string) error
	// Reads group chat messagses, on or after startTime and on or before endTime, from the database.
	//
	// May not return all messages. If the returned `bool` is true, there may be
	// messages remaining (set `startTime` to the latest `message.Timestamp` and try again).
	ReadMessages(GroupID, startTime UnixMillis, endTime UnixMillis) ([]Message, bool, error)
	// Creates a new poll in the group, replacing the old one (if any).
	//
	// Returns an error if the operation could not be completed.
	CreatePoll(GroupID, Poll) error
	// Deletes the active poll in a group, if one exists.
	//
	// Returns an error if the operation could not be completed.
	DeletePoll(GroupID) error
	// Creates an availability for a user in a group.
	CreateAvailability(GroupID /*, Activity*/) error
	// Creates an activity in a group.
	CreateActivity(GroupID /*, Activity*/) error
	// Deletes a group from the database, if it exists.
	//
	// Returns an error if the operation could not be completed.
	DeleteGroup(GroupID) error
	// Creates a new chat message in the group.
	//
	// Returns an error if the `message.GroupID` and `message.Timestamp` are not
	// unique, or if the operation could not be completed.
	CreateMessage(Message) error
}

// An AWS non-volatile database service.
type DynamoDB struct {
	groups   dynamo.Table
	users    dynamo.Table
	messages dynamo.Table
}

type localConfigProvider struct{}

func (p localConfigProvider) ClientConfig(serviceName string, cfgs ...*aws.Config) client.Config {

	return client.Config{Endpoint: "http://localhost:8000", Config: aws.NewConfig()}
}

// Passing a `nil` session means use DynamoDB local (default port).
func NewDynamoDB(sess *session.Session) *DynamoDB {
	var db *dynamo.DB
	if sess == nil {
		// localConfigProvider(struct{}{})
		endpoint := "http://localhost:8000"
		sess, err := session.NewSession(
			&aws.Config{Region: aws.String(GetRegion()), Endpoint: &endpoint, Credentials: credentials.NewCredentials(&credentials.StaticProvider{credentials.Value{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
				ProviderName:    "dummy",
			}})},
		)
		if err != nil {
			log.Fatal(err)
		}
		db = dynamo.New(sess)
		if err := db.CreateTable("Groups", Group{}).Run(); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable("Users", User{}).Run(); err != nil {
			log.Fatal(err)
		}
		if err := db.CreateTable("Messages", Message{}).Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		db = dynamo.New(sess, &aws.Config{Region: aws.String(GetRegion())})
	}

	return &DynamoDB{
		groups:   db.Table("Groups"),
		users:    db.Table("Users"),
		messages: db.Table("Messages"),
	}
}

// Creates a new user in the database.
//
// Returns an error if a user with the same `UserID`
// already exists, or if the operation may have failed.
func (dynamoDB *DynamoDB) CreateUser(userInfo User) error {
	return dynamoDB.users.Put(userInfo).If("attribute_not_exists(UserID)").Run()
}

// Reads a user from the database.
//
// Returns a nil `*User` if no such user exists. Returns
// an error if the operation could not be completed.
func (dynamoDB *DynamoDB) ReadUser(userID UserID) (*User, error) {
	var user User
	err := dynamoDB.users.Get("UserID", userID).One(&user)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, nil
	}
	return &user, err
}

// Deletes a user from the database, if it exists.
//
// Returns an error if the operation could not be completed.
func (dynamoDB *DynamoDB) DeleteUser(userID UserID) error {
	return dynamoDB.users.Delete("UserID", userID).If("attribute_exists(UserID)").Run()
}

// Creates a new group in the database.
//
// Returns an error if a group with the same `GroupIP`
// already exists, or if the operation may have failed.
func (dynamoDB *DynamoDB) CreateGroup(groupInfo Group) error {
	return dynamoDB.groups.Put(groupInfo).If("attribute_not_exists(GroupID)").Run()
}

// Reads a group from the database.
//
// Returns a nil `*Group` if no such group exists. Returns
// an error if the operation could not be completed.
func (dynamoDB *DynamoDB) ReadGroup(groupID GroupID) (*Group, error) {
	var group Group
	err := dynamoDB.groups.Get("GroupID", groupID).One(&group)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, nil
	}
	return &group, err
}

func (dynamoDB *DynamoDB) UpdateGroupName(groupID GroupID, name string) error {
	panic("unimplemented")
}

func (dynamoDB *DynamoDB) ReadMessages(groupID GroupID, startTime UnixMillis, endTime UnixMillis) ([]Message, bool, error) {
	var messages []Message
	const limit = 5
	err := dynamoDB.messages.Get("GroupID", groupID).Range("Timestamp", "BETWEEN", startTime, endTime).Limit(limit).All(&messages)
	return messages, len(messages) >= limit, err
}

// Deletes a group from the database, if it exists.
//
// Returns an error if the operation could not be completed.
func (dynamoDB *DynamoDB) DeleteGroup(groupID GroupID) error {
	err := dynamoDB.groups.Delete("GroupID", groupID).If("attribute_exists(UserID)").Run()
	return err
}

// Inserts a new schedule to the database, if it does not exist
// Returns an error if schedule already exists in database
func (dynamoDB *DynamoDB) InsertNewSchedule(userID UserID, scheduleInfo Schedule) error {
	err := dynamoDB.users.Update("UserID", userID).Append("Schedules", scheduleInfo).Run()
	return err
}

// Updates the group information in the database, if it exist
// Returns an error if group does not exist in database
func (dynamoDB *DynamoDB) UpdateGroupInfo(groupID GroupID, newInfo Group) error {
	err := dynamoDB.groups.Update("GroupID", groupID).Set("Group", newInfo).Run()
	return err
}

// Updates the user information in the database, if it exist
// Returns an error if user does not exist in database
func (dynamoDB *DynamoDB) UpdateUserInfo(userID UserID, newInfo User) error {
	err := dynamoDB.users.Update("UserID", userID).Set("User", newInfo).Run()
	return err
}

// Deletes a user from the group database, if the user exists in that group.
//
// Returns an error if the operation could not be completed.
func (dynamoDB *DynamoDB) DeleteUserFromGroup(userInfo User, groupID GroupID) error {
	//Check if group exists, check if user exists
	err := dynamoDB.groups.Update("GroupID", groupID).DeleteFromSet("Users", userInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) CreateActivity(groupID GroupID /*, activity Activity*/) error {
	return fmt.Errorf("unimplemented")
}

func (dynamoDB *DynamoDB) CreateAvailability(groupID GroupID /*, availability Availability*/) error {
	return fmt.Errorf("unimplemented")
}

func (dynamoDB *DynamoDB) CreatePoll(groupID GroupID, poll Poll) error {
	return fmt.Errorf("unimplemented")
}

func (dynamoDB *DynamoDB) DeletePoll(groupID GroupID) error {
	return fmt.Errorf("unimplemented")
}

func (dynamoDB *DynamoDB) CreateMessage(message Message) error {
	return fmt.Errorf("unimplemented")
}

func printDatabase(database Database) error {
	out, err := json.Marshal(database)
	fmt.Println(string(out))
	return err
}

// An in-memory volatile database.
type MemoryDatabase struct {
	users    map[UserID]User
	groups   map[GroupID]Group
	messages map[memoryMessageID]Message
	mu       sync.Mutex
}

type memoryMessageID struct {
	GroupID   GroupID
	Timestamp uint64
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:  make(map[UserID]User),
		groups: make(map[GroupID]Group),
	}
}

func (memoryDatabase *MemoryDatabase) CreateUser(user User) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	if _, ok := memoryDatabase.users[user.UserID]; ok {
		return fmt.Errorf("user already exists")
	}
	memoryDatabase.users[user.UserID] = user
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadUser(userID UserID) (*User, error) {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	user, ok := memoryDatabase.users[userID]
	if ok {
		return &user, nil
	} else {
		return nil, nil
	}
}

func (memoryDatabase *MemoryDatabase) DeleteUser(userID UserID) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	delete(memoryDatabase.users, userID)
	return nil
}

func (memoryDatabase *MemoryDatabase) CreateGroup(group Group) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	if _, ok := memoryDatabase.groups[group.GroupID]; ok {
		return fmt.Errorf("group already exists")
	}
	memoryDatabase.groups[group.GroupID] = group
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadGroup(groupId GroupID) (*Group, error) {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	group, ok := memoryDatabase.groups[groupId]
	if ok {
		return &group, nil
	} else {
		return nil, nil
	}
}

func (memoryDatabase *MemoryDatabase) UpdateGroupName(groupID GroupID, name string) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	group, ok := memoryDatabase.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	group.Name = name
	memoryDatabase.groups[groupID] = group
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadMessages(groupID GroupID, startTime UnixMillis, endTime UnixMillis) ([]Message, bool, error) {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	var messages []Message
	var more bool
	// Okay to do inefficient linear table scan on mock database.
	for _, message := range memoryDatabase.messages {
		if message.GroupID != groupID || message.Timestamp < startTime || message.Timestamp > endTime {
			continue
		}
		if len(messages) >= 5 {
			more = true
			break
		}
		messages = append(messages, message)
	}
	return messages, more, nil
}

func (memoryDatabase *MemoryDatabase) CreateActivity(groupID GroupID /*, activity Activity*/) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	// TODO: unimplemented.
	return nil
}

func (memoryDatabase *MemoryDatabase) CreateAvailability(groupID GroupID /*, availability Availability*/) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	_, ok := memoryDatabase.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	//group.Availabilities = append(group.Availabilities, availability)
	return nil
}

func (memoryDatabase *MemoryDatabase) CreatePoll(groupID GroupID, poll Poll) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	group, ok := memoryDatabase.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	group.Poll = &poll
	memoryDatabase.groups[groupID] = group
	return nil
}

func (memoryDatabase *MemoryDatabase) DeletePoll(groupID GroupID) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	group, ok := memoryDatabase.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	group.Poll = nil
	memoryDatabase.groups[groupID] = group
	return nil
}

func (memoryDatabase *MemoryDatabase) DeleteGroup(groupID GroupID) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	delete(memoryDatabase.groups, groupID)
	return nil
}

func (memoryDatabase *MemoryDatabase) CreateMessage(message Message) error {
	id := memoryMessageID{
		GroupID:   message.GroupID,
		Timestamp: message.Timestamp,
	}
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	if _, ok := memoryDatabase.messages[id]; ok {
		return fmt.Errorf("message already exists")
	}
	memoryDatabase.messages[id] = message
	return nil
}

func GetRegion() string {
	var region = os.Getenv("AWS_REGION")

	if region == "" {
		return "us-east-1"
	}

	return region
}
