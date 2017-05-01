package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/NZguy/in-class/tasksvr/models/tasks"
)

//HandleTasks will handle requests for the /v1/tasks resource
func (ctx *Context) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		newtask := &tasks.NewTask{}
		if err := decoder.Decode(newtask); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := newtask.Validate(); err != nil {
			http.Error(w, "error validating task: "+err.Error(), http.StatusBadRequest)
			return
		}

		task, err := ctx.TasksStore.Insert(newtask)
		if err != nil {
			http.Error(w, "error inserting task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(task) // encode back to client

	case "GET":
		tasks, err := ctx.TasksStore.GetAll()
		if err != nil {
			http.Error(w, "error getting tasks", http.StatusInternalServerError)
		}
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(tasks)
	}
}

//HandleSpecificTask will handle requests for the /v1/tasks/some-task-id resource
func (ctx *Context) HandleSpecificTask(w http.ResponseWriter, r *http.Request) {
	_, id := path.Split(r.URL.Path)

	switch r.Method {
	case "GET":
		task, err := ctx.TasksStore.Get(id)
		if err != nil {
			http.Error(w, "error finding task: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Should make a function for this as its super common to write json
		// Question: Should go in the handlers package in its own file with all other constants
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(task)
	case "PATCH":
		decoder := json.NewDecoder(r.Body)
		task := &tasks.Task{}
		if err := decoder.Decode(task); err != nil {
			http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		task.ID = id

		if err := ctx.TasksStore.Update(task); err != nil {
			http.Error(w, "error updating: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
