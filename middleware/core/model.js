var path = require('path');
var Sequelize = require('sequelize');
var _ = require('underscore');

var database = require('../../config/database');

global.Model = Sequelize;

global.model = new Sequelize(
    database.database,
    database.username,
    database.password,
    _.extend(database, {
        logging: process.env['NODE_ENV'] !== 'production' && console.log
    })
);

require("fs").readdirSync("./models").forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require("../models/" + file);
});
