package db

import {
	"fmt"
	"reflect"
}

type Tables struct {
	groups	Table
	users	Table

}

func initializeTables(db reflect.Type) {
	Tables.groups := db.Table("Groups")
	Tables.users := db.Table("Users")
}	

func insertNewUser(user userInfo) {
	err := Tables.users.Put(userInfo).Run()
}

func insertNewGroup(group groupInfo) {
	err := Tables.groups.Put(groupInfo).Run()
}

func insertNewSchedule(string userID, schedule scheduleInfo) {
	err = Tables.users.Update("UserID", userID)
					  .Append("Schedules", scheduleInfo)
					  .Run()
} 

func updateUserInfo(string userID, map[string]string newInfo) {

}

func updateGroup() {

}

func deleteUserFromGroup(user userInfo, group groupInfo) {

} 

func deleteGroup(group groupIndo) {

}