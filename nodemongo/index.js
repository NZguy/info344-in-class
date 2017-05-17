"use strict";

const express = require('express');
const morgan = require('morgan');
const cors = require('cors');
const bodyParser = require('body-parser');
const mongodb = require('mongodb');

// We could use a different database by changing the file we import from
const TaskStore = require('./models/tasks/mongostore.js');

const port = process.env.PORT || 80;
const host = process.env.HOST || '';
const mongoAddr = process.env.MONGOADDR || 'localhost:27017';

//create an Experss application
const app = express();
//add request logging
app.use(morgan(process.env.LOGFORMAT || 'dev'));
//add CORS headers
app.use(cors());
//add middleware that parses
//any JSON posted to this app.
//the parsed data will be available
//on the req.body property
app.use(bodyParser.json());

//TODO: connect to the Mongo database
//add the tasks handlers
//and start listening for HTTP requests
mongodb.MongoClient.connect(`mongodb://${mongoAddr}/demo`)
    .then(db => {
        let colTasks = db.collection('tasks');
        let store = new TaskStore(colTasks);
        let handlers = require('./handlers/tasks.js');
        // Router is added to application
        // First parameter is a prefix for all paths in returned router
        app.use('/api', handlers(store));

        // Default error handler sends stack traces to users
        // If we pass a function wiht 4 parameters to app.use, express assumes its a error handler
        app.use((err, req, res, next) => {
            console.error(err);
            // If an error status exists, ues it instead
            res.status(err.status || 500).send(err.message);
        });

        app.listen(port, host, () => {
            console.log(`server is listening at http://${host}:${port}...`);
        })
    })
    .catch(err => {
        console.error(err);
    });

    // Use mocha and maybe chai for testing

