package main

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/dynamo"
)

type ID = uint64

type Database struct {
	groups dynamo.Table
	users  dynamo.Table
}

func initializeNewDatabase(db *dynamo.DB) *Database {
	return &Database{
		groups: db.Table("Groups"),
		users:  db.Table("Users"),
	}
}

func (database *Database) InsertNewUser(userInfo User) error {
	err := database.users.Put(userInfo).Run()
	return err
}

func (database *Database) InsertNewGroup(groupInfo Group) error {
	err := database.groups.Put(groupInfo).Run()
	return err
}

func (database *Database) InsertNewSchedule(userID ID, scheduleInfo Schedule) error {
	err := database.users.Update("UserID", userID).Append("Schedules", scheduleInfo).Run()
	return err
}

func (database *Database) UpdateGroupInfo(groupID ID, newInfo Group) error {
	err := database.groups.Update("GroupID", groupID).Set("Group", newInfo).Run()
	return err
}

func (database *Database) UpdateUserInfo(userID ID, newInfo User) error {
	err := database.users.Update("UserID", userID).Set("User", newInfo).Run()
	return err
}

func (database *Database) deleteUserFromGroup(userInfo User, groupID ID) error {
	//Check if group exists, check if user exists
	err := database.groups.Update("GroupID", groupID).DeleteFromSet("Users", userInfo).Run()
	return err
}

func (database *Database) deleteGroup(groupInfo Group) error {
	err := database.groups.Delete("GroupID", groupInfo.GroupID).Run()
	return err
}

func printDatabase(database Database) error {
	out, err := json.Marshal(database)
	fmt.Println(string(out))
	return err
}
