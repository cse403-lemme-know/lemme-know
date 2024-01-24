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

type Database interface {
	CreateUser(User) error
	ReadUser(UserID) (*User, error)
	DeleteUser(UserID) error
	CreateGroup(Group) error
	ReadGroup(GroupID) (*Group, error)
	DeleteGroup(GroupID) error
}

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

func (dynamoDB *DynamoDB) InsertNewUser(userInfo User) error {
	err := dynamoDB.users.Put(userInfo).Run()
	return err
}

func (dynamoDB *DynamoDB) InsertNewGroup(groupInfo Group) error {
	err := dynamoDB.groups.Put(groupInfo).Run()
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

func (dynamoDB *DynamoDB) deleteGroup(groupInfo Group) error {
	err := dynamoDB.groups.Delete("GroupID", groupInfo.GroupID).Run()
	return err
}

func printDatabase(database Database) error {
	out, err := json.Marshal(database)
	fmt.Println(string(out))
	return err
}
