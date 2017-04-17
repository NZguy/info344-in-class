package tasks

import (
	"fmt"
	"time"
)

//NewTask represents a new task posted to the server
// Struct fields must be exported (capitalized) to be encoded to json
type NewTask struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"` // Allows us to use lowercase names in json
}

//Task represents a task stored in the database
type Task struct {
	ID         interface{} `json:"id" bson:"_id"` // interface is go's any type
	Title      string      `json:"title"`
	Tags       []string    `json:"tags"`
	CreatedAt  time.Time   `json:"createdAt"`
	ModifiedAt time.Time   `json:"modifiedAt"`
	Complete   bool        `json:"complete"`
}

//Validate will validate the NewTask
func (nt *NewTask) Validate() error {
	//Title fiend must be non-zero length
	if len(nt.Title) == 0 {
		return fmt.Errorf("title must be something")
	}
	return nil
}

//ToTask converts a NewTask to a Task
func (nt *NewTask) ToTask() *Task {
	task := &Task{ // Shorthand to create new task and gives us a reference to it
		Title:      nt.Title,
		Tags:       nt.Tags,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	return task
}
