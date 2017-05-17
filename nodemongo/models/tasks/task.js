"use strict";

//TODO: implement a Task class
//for the task model, and export
//it from this module

class Task {
    constructor(props) {
        // Take props are make them the fields of this object
        Object.assign(this, props);
    }

    validate(){
        if(!this.title){
            return new Error('You must supply a title');
        }
    }
}

// Allows us to import task from other files
module.exports = Task;
