var express = require('express');
var bodyParser = require('body-parser')

var app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

require('./core/models');

app.use('/:domain/:secret',
    require('./core/actions')
);

app.listen(9001, function () {
    console.log('Api server is starting now.')
});