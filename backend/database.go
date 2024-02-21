package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

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
	// Transactionally updates a user.
	//
	// Returns an error if the operation could not be completed.
	UpdateUser(UserID, func(user *User) error) error
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
	// Transactionally updates a group.
	//
	// Returns an error if the operation could not be completed.
	UpdateGroup(GroupID, func(group *Group) error) error
	// Reads group chat messagses, on or after startTime and on or before endTime, from the database.
	//
	// May not return all messages. If the returned `bool` is true, there may be
	// messages remaining (set `startTime` to the latest `message.Timestamp` and try again).
	ReadMessages(GroupID, startTime UnixMillis, endTime UnixMillis) ([]Message, bool, error)
	// Returns an error if the operation could not be completed.
	DeleteGroup(GroupID) error
	// Creates a new chat message in the group.
	//
	// Returns an error if the `message.GroupID` and `message.Timestamp` are not
	// unique, or if the operation could not be completed.
	CreateMessage(Message) error
	// Reads the UserID associated with a connection, or returns nil if none exists.
	//
	// Returns an error if the operation could not be completed.
	ReadConnection(ConnectionID) (*UserID, error)
	// Writes or overwrites the UserID associated with a connection.
	//
	// Returns an error if the operation could not be completed.
	WriteConnection(ConnectionID, UserID) error
	// Deletes a connection.
	//
	// Returns an error if the operation could not be completed.
	DeleteConnection(ConnectionID) error
	// Reads the value of a variable (possibly empty string if empty or nonexistent).
	ReadVariable(string) (string, error)
	// Overwrites the value of a variable.
	WriteVariable(string, string) error
}

// An AWS non-volatile database service.
type DynamoDB struct {
	groups      dynamo.Table
	users       dynamo.Table
	messages    dynamo.Table
	connections dynamo.Table
	variables   dynamo.Table
}

const groupTableName = "lemmeknow-groups"
const userTableName = "lemmeknow-users"
const messageTableName = "lemmeknow-messages"
const connectionTableName = "lemmeknow-connections"
const variableTableName = "lemmeknow-variables"

// Passing a `nil` session means use DynamoDB local (default port).
func NewDynamoDB(sess *session.Session) *DynamoDB {
	var db *dynamo.DB
	if sess == nil {
		endpoint := "http://localhost:8000"
		sess, err := session.NewSession(
			&aws.Config{Region: aws.String(GetRegion()), Endpoint: &endpoint, Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
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

		// Ingnore errors (e.g. duplicate table)
		_ = db.CreateTable(groupTableName, Group{}).Run()
		_ = db.CreateTable(userTableName, User{}).Run()
		_ = db.CreateTable(messageTableName, Message{}).Run()
		_ = db.CreateTable(connectionTableName, Connection{}).Run()
		_ = db.CreateTable(variableTableName, Variable{}).Run()
	} else {
		db = dynamo.New(sess, &aws.Config{Region: aws.String(GetRegion())})
	}

	return &DynamoDB{
		groups:      db.Table(groupTableName),
		users:       db.Table(userTableName),
		messages:    db.Table(messageTableName),
		connections: db.Table(connectionTableName),
		variables:   db.Table(variableTableName),
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

func (dynamoDB *DynamoDB) UpdateUser(userID UserID, transaction func(*User) error) error {
	governor := 0
	for {
		user, err := dynamoDB.ReadUser(userID)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user not found")
		}

		oldCount := user.UpdateCount
		if err := transaction(user); err != nil {
			return err
		}
		user.UpdateCount = oldCount + 1

		err = dynamoDB.users.Put(user).If("UpdateCount = ?", oldCount).Run()
		if err != nil && dynamo.IsCondCheckFailed(err) {
			// Retry the transaction.
			governor += 1
			if governor > 16 {
				return fmt.Errorf("too many retries")
			}
			continue
		}
		return err
	}
}

// Deletes a user from the database, if it exists.
//
// Returns an error if the operation could not be completed.
func (dynamoDB *DynamoDB) DeleteUser(userID UserID) error {
	return dynamoDB.users.Delete("UserID", userID).If("attribute_exists(UserID)").Run()
}

// Creates a new group in the database.
//
// Returns an error if a group with the same `GroupID`
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

func (dynamoDB *DynamoDB) UpdateGroup(groupID GroupID, transaction func(*Group) error) error {
	governor := 0
	for {
		group, err := dynamoDB.ReadGroup(groupID)
		if err != nil {
			return err
		}
		if group == nil {
			return fmt.Errorf("group not found")
		}

		oldCount := group.UpdateCount
		if err := transaction(group); err != nil {
			return err
		}
		group.UpdateCount = oldCount + 1

		err = dynamoDB.groups.Put(group).If("UpdateCount = ?", oldCount).Run()
		if err != nil && dynamo.IsCondCheckFailed(err) {
			// Retry the transaction.
			governor += 1
			if governor > 16 {
				return fmt.Errorf("too many retries")
			}
			continue
		}
		return err
	}
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
	err := dynamoDB.groups.Delete("GroupID", groupID).If("attribute_exists(GroupID)").Run()
	return err
}

func (dynamoDB *DynamoDB) CreateMessage(message Message) error {
	return dynamoDB.messages.Put(message).If("attribute_not_exists($)", "Timestamp").Run()
}

func (dynamoDB *DynamoDB) ReadConnection(connectionID ConnectionID) (*UserID, error) {
	var connection Connection
	err := dynamoDB.connections.Get("ConnectionID", connectionID).One(&connection)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, nil
	}
	return &connection.UserID, err
}

func (dynamoDB *DynamoDB) WriteConnection(connectionID ConnectionID, userID UserID) error {
	connection := Connection{
		ConnectionID: connectionID,
		UserID:       userID,
	}
	return dynamoDB.connections.Put(connection).If("attribute_not_exists(ConnectionID)").Run()
}

func (dynamoDB *DynamoDB) DeleteConnection(connectionID ConnectionID) error {
	err := dynamoDB.connections.Delete("ConnectionID", connectionID).If("attribute_exists(ConnectionID)").Run()
	return err
}

func (dynamoDB *DynamoDB) ReadVariable(name string) (string, error) {
	var variable Variable
	err := dynamoDB.variables.Get("Name", name).One(&variable)

	if errors.Is(err, dynamo.ErrNotFound) {
		return "", nil
	}
	return variable.Value, err
}

func (dynamoDB *DynamoDB) WriteVariable(name string, value string) error {
	variable := Variable{
		Name:  name,
		Value: value,
	}
	return dynamoDB.variables.Put(variable).If("attribute_not_exists($)", "Name").Run()
}

func printDatabase(database Database) error {
	out, err := json.Marshal(database)
	fmt.Println(string(out))
	return err
}

// An in-memory volatile database.
type MemoryDatabase struct {
	users       map[UserID]User
	groups      map[GroupID]Group
	messages    map[memoryMessageID]Message
	connections map[ConnectionID]UserID
	variables   map[string]string
	mu          sync.Mutex
}

type memoryMessageID struct {
	GroupID   GroupID
	Timestamp uint64
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:       make(map[UserID]User),
		groups:      make(map[GroupID]Group),
		messages:    make(map[memoryMessageID]Message),
		connections: make(map[ConnectionID]UserID),
		variables:   make(map[string]string),
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

func (memoryDatabase *MemoryDatabase) UpdateUser(userID UserID, transaction func(*User) error) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	user, ok := memoryDatabase.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	if err := transaction(&user); err != nil {
		return err
	}
	memoryDatabase.users[userID] = user
	return nil
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

func (memoryDatabase *MemoryDatabase) UpdateGroup(groupID GroupID, transaction func(*Group) error) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	group, ok := memoryDatabase.groups[groupID]
	if !ok {
		return fmt.Errorf("group not found")
	}
	if err := transaction(&group); err != nil {
		return err
	}
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

func (memoryDatabase *MemoryDatabase) ReadConnection(connectionID ConnectionID) (*UserID, error) {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	userID, ok := memoryDatabase.connections[connectionID]
	if ok {
		return &userID, nil
	} else {
		return nil, nil
	}
}

func (memoryDatabase *MemoryDatabase) WriteConnection(connectionID ConnectionID, userID UserID) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	memoryDatabase.connections[connectionID] = userID
	return nil
}

func (memoryDatabase *MemoryDatabase) DeleteConnection(connectionID ConnectionID) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	delete(memoryDatabase.connections, connectionID)
	return nil
}

func (memoryDatabase *MemoryDatabase) ReadVariable(name string) (string, error) {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	value, ok := memoryDatabase.variables[name]
	if ok {
		return value, nil
	} else {
		return "", nil
	}
}

func (memoryDatabase *MemoryDatabase) WriteVariable(name string, value string) error {
	memoryDatabase.mu.Lock()
	defer memoryDatabase.mu.Unlock()
	memoryDatabase.variables[name] = value
	return nil
}

func GetRegion() string {
	var region = os.Getenv("AWS_REGION")

	if region == "" {
		return "us-east-1"
	}

	return region
}

// Generates a positive ID number that won't experience precision
// issues for the client's JavaScript number representation.
func GenerateID() uint64 {
	const maxSafeInteger = 9007199254740991
	return uint64(rand.Int63n(maxSafeInteger + 1))
}
