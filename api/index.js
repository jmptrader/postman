var express = require('express');
var bodyParser = require('body-parser');
var port = require('../config/domain').api.port;

var app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

require('./core/models');
require('./core/command');

app.use('/:domain/:secret',
    require('./core/actions')
);

app.listen(port, function () {
    console.log('Api server is starting now.')
});