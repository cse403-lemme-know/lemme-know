package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type UserID = uint64
type GroupID = uint64

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
	// Returns an error if a group with the same `GroupIP`
	// already exists, or if the operation may have failed.
	CreateGroup(Group) error
	// Reads a group from the database.
	//
	// Returns a nil `*Group` if no such group exists. Returns
	// an error if the operation could not be completed.
	ReadGroup(GroupID) (*Group, error)
	// Deletes a group from the database, if it exists.
	//
	// Returns an error if the operation could not be completed.
	DeleteGroup(GroupID) error
}

// An AWS non-volatile database service.
type DynamoDB struct {
	groups dynamo.Table
	users  dynamo.Table
}

func NewDynamoDB(sess *session.Session) Database {
	// TODO: Hard-coded region.
	db := dynamo.New(sess, &aws.Config{Region: aws.String("us-east-1")})
	return &DynamoDB{
		groups: db.Table("Groups"),
		users:  db.Table("Users"),
	}
}

func (dynamoDB *DynamoDB) CreateUser(userInfo User) error {
	err := dynamoDB.users.Put(userInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) ReadUser(userId UserID) (*User, error) {
	// TODO: unimplemented.
	return nil, nil
}

func (dynamoDB *DynamoDB) DeleteUser(userID UserID) error {
	err := dynamoDB.users.Delete("UserID", userID).Run()
	return err
}

func (dynamoDB *DynamoDB) CreateGroup(groupInfo Group) error {
	err := dynamoDB.groups.Put(groupInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) ReadGroup(groupID GroupID) (*Group, error) {
	// TODO: unimplemented.
	return nil, nil
}

func (dynamoDB *DynamoDB) DeleteGroup(groupID GroupID) error {
	err := dynamoDB.groups.Delete("GroupID", groupID).Run()
	return err
}

func (dynamoDB *DynamoDB) InsertNewSchedule(userID UserID, scheduleInfo Schedule) error {
	err := dynamoDB.users.Update("UserID", userID).Append("Schedules", scheduleInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) UpdateGroupInfo(groupID GroupID, newInfo Group) error {
	err := dynamoDB.groups.Update("GroupID", groupID).Set("Group", newInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) UpdateUserInfo(userID UserID, newInfo User) error {
	err := dynamoDB.users.Update("UserID", userID).Set("User", newInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) deleteUserFromGroup(userInfo User, groupID GroupID) error {
	//Check if group exists, check if user exists
	err := dynamoDB.groups.Update("GroupID", groupID).DeleteFromSet("Users", userInfo).Run()
	return err
}

func printDatabase(database Database) error {
	out, err := json.Marshal(database)
	fmt.Println(string(out))
	return err
}

// An in-memory volatile database.
type MemoryDatabase struct {
	users  map[UserID]User
	groups map[GroupID]Group
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:  make(map[UserID]User),
		groups: make(map[GroupID]Group),
	}
}

func (memoryDatabase *MemoryDatabase) CreateUser(user User) error {
	if _, ok := memoryDatabase.users[user.UserID]; ok {
		return fmt.Errorf("user already exists")
	}
	memoryDatabase.users[user.UserID] = user
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadUser(userID UserID) (*User, error) {
	user, ok := memoryDatabase.users[userID]
	if ok {
		return &user, nil
	} else {
		return nil, nil
	}
}

func (memoryDatabase *MemoryDatabase) DeleteUser(userID UserID) error {
	delete(memoryDatabase.users, userID)
	return nil
}

func (memoryDatabase *MemoryDatabase) CreateGroup(group Group) error {
	if _, ok := memoryDatabase.groups[group.GroupID]; ok {
		return fmt.Errorf("group already exists")
	}
	memoryDatabase.groups[group.GroupID] = group
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadGroup(groupId GroupID) (*Group, error) {
	group, ok := memoryDatabase.groups[groupId]
	if ok {
		return &group, nil
	} else {
		return nil, nil
	}
}

func (memoryDatabase *MemoryDatabase) DeleteGroup(groupID GroupID) error {
	delete(memoryDatabase.groups, groupID)
	return nil
}
