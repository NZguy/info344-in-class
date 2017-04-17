package tasks

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
	task := &Task{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(ID).One(task) // we have to tell mongo what struct we want it to fill out
	return task, err
}
