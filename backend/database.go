package main

import (
	"github.com/guregu/dynamo"
)

type UserID = uint64

type Database struct {
	groups	dynamo.Table
	users	dynamo.Table
}

func initializeNewDatabase(db *dynamo.DB) *Database{
	return &Database{
		groups: db.Table("Groups"),
		users: db.Table("Users"),
	}
}	

func (database *Database) InsertNewUser(userInfo User) error {
	return database.users.Put(userInfo).Run()
}

func (database *Database) InsertNewGroup(groupInfo Group) error {
	return database.groups.Put(groupInfo).Run()
}

func (database *Database) InsertNewSchedule(userID UserID, scheduleInfo Schedule) error {
	return database.users.Update("UserID", userID).Append("Schedules", scheduleInfo).Run()
} 

func (database *Database) UpdateUserInfo(userID UserID, newInfo map[string]string) error {
	return nil
}

func (database *Database) updateGroup() error {
	return nil
}

func (database *Database) deleteUserFromGroup(userInfo User, groupInfo Group) error {
	return nil
} 

func (database *Database) deleteGroup(groupInfo Group) error {
	return nil
}