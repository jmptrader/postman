var express = require('express');
var app = express();

require('./core/models');

app.use('/:domain/:secret',
    require('./core/actions')
);

app.listen(9001, function () {
    console.log('Api server is starting now.')
});