"use strict";

const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const zips = require('./zips.json');
console.log('loaded %d zips', zips.length);

const zipCityIndex = zips.reduce((accumulator, record) => {
    let cityLower = record.city.toLowerCase();
    // accumulator is a "global" variable that we put all the data into
    let zipsForCity = accumulator[cityLower];
    if (!zipsForCity){
        // Cascading assignment zipsForCity becomes [], then index[cityLower] becomes zipsForCity
        accumulator[cityLower] = zipsForCity = [];
    }
    zipsForCity.push(record);
    return accumulator;
}, {});
console.log('there are %d zips in Seattle', zipCityIndex.seattle.length);

const app = express();

// In js we can use the result of an or as the result. So if port is true than it will be used.
// If port is false, 80 will be used.
// Called null coalesce
const port = process.env.PORT || 80;
const host = process.env.HOST || '';

// .use applies the given function to all requests
app.use(morgan('dev')); // morgan returns a middleware function to be used
app.use(cors());

app.get('/zips/city/:cityName', (req, res) => {
    let zipsForCity = zipCityIndex[req.params.cityName.toLowerCase()];
    if (!zipsForCity){
        res.status(404).send('invalid city name');
    } else {
        // Converts response to json and automatically sets content type header
        res.json(zipsForCity);
    }
});

// : is a wildcard, the string entered will be availible as the specified variable name
// app.get only responds to get requests
app.get('/hello/:name', (req, res) => {
    res.send(`Hello, ${req.params.name}!`);
});

app.listen(port, host, () => {
    // JS has template strings now, use this ``
    console.log(`Server is listening at http://${host}:${port}`);
});
