"use strict";

const mongodb = require('mongodb'); //for mongodb.ObjectID()

/**
 * MongoStore is a concrete store for Task models
 */
class MongoStore {
    /**
     * Constructs a new MongoStore
     * @param {mongodb.Collection} collection
     */
    constructor(collection) {
        this.collection = collection;
    }

    /**
     * getAll returns all tasks in the store
     */
    getAll() {
        return this.collection.find().toArray();
    }

    /**
     * insert inserts a new Task into the store
     * @param {Task} task
     */
    insert(task) {
        return this.collection.insert(task);
    }

    /**
     * setComplete sets the complete status of the task
     * @param {string} id
     * @param {bool} complete
     */
    async setComplete(id, complete) {
        // async tag must be used in any function that uses await

        // Prevents mongo from returning non updated version of document
        let options = {returnOriginal: false}
        let updates = {$set: {complete: complete}};
        // Q: Completions?
        let oid = new mongodb.ObjectID(id);
        // If there is an error, await will throw a js exception
        let result = await this.collection.findOneAndUpdate({_id: oid}, updates, options);
        return result.value; // Updated document
    }

    /**
     * delete deletes the task with the given object ID
     * @param {string} id
     */
    delete(id) {
        return this.collection.deleteOne({_id: new mongodb.ObjectID(id)});
    }
}

//export the class
module.exports = MongoStore;
