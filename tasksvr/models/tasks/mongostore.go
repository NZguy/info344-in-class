package tasks

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo sucks as many to many relations, sql actually has JOINS

// Implementation of store interface is assumed by the structure of the file

type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
}

func (ms *MongoStore) Insert(newtask *NewTask) (*Task, error) {
	t := newtask.ToTask()
	t.ID = bson.NewObjectId()
	// Insert to collection in database
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(t)
	return t, err
}

func (ms *MongoStore) Get(ID interface{}) (*Task, error) {
	// Mongo needs the string ID converted to its version of an object ID
	// If ID is equal to a string, ok will return true, and sID will be equal to a string version of ID
	if sID, ok := ID.(string); ok {
		ID = bson.ObjectIdHex(sID) // Mongo does weird shit with interfaces, it converts interfaces to hex and they must be converted back
	}
	task := &Task{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task) // we have to tell mongo what struct we want it to fill out
	return task, err
}

func (ms *MongoStore) GetAll() ([]*Task, error) {
	tasks := []*Task{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&tasks) // Find all needs a pointer becuase its going to edit tasks to populate
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ms *MongoStore) Update(task *Task) error {
	if sID, ok := task.ID.(string); ok { // Type assertion
		task.ID = bson.ObjectIdHex(sID)
	}
	task.ModifiedAt = time.Now()
	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName)
	// bson.M creates a bson map
	// we only update the fields tat we want to be updatable
	updates := bson.M{"$set": bson.M{"complete": task.Complete, "modifiedat": task.ModifiedAt}}
	return col.UpdateId(task.ID, updates) // returning becuase will return error if one exists
}
