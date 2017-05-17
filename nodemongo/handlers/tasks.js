"use strict";

const express = require('express');
// Q: You got completions here, why don't I?
const Task = require('../models/tasks/task.js');

//export a function from this module
//that accepts a tasks store implementation
module.exports = function(store) {
    //TODO: create an express.Router()
    //add handlers to it
    //and return the configured Router
    let router = express.Router();

    router.get('/v1/tasks', async (req, res, next) => {
        // Handle the error that may occur in mongostore
        try {
            // throw an error with more details
            //throw {status: 550, message: 'testing'};
            let tasks = await store.getAll();
            res.json(tasks);
        } catch(err) {
            // Express uses next function to do centralized error handling
            next(err);
        }

        // Without await you would need to use promise syntax,
        // this is syntax sugar, the code would be functionally the same
        // store.getAll()
        //     .then(tasks => {
        //         res.json(tasks);
        //     })
        //     .catch(next);
    });
                                                // Inline anon function big arrow
    router.post('/v1/tasks', async (req, res, next) => {
        try {
            // req.body is set by the bodyparser middleware
            let task = new Task(req.body);
            let err = task.validate();
            if (err) {
                res.status(400).send(err.message);
            }else{
                let result = await store.insert(task);
                res.json(task);
            }
        } catch(err) {
            next(err);
        }
    });

    // :syntax is called a router wildcard
    router.patch('/v1/tasks/:taskID', async (req, res, next) => {
        let taskID = req.params.taskID;
        try {
            let result = await store.setComplete(taskID, req.body.complete);
            res.json(result);
        } catch(err) {
            next(err);
        }
    });

    router.delete('/v1/tasks/:taskID', async (req, res, next) => {
        let taskID = req.params.taskID;
        try {
             await store.delete(taskID);
             res.send(`deleted task ${taskID}`);
        } catch(err){
            next(err);
        }
    });

    return router;
};
