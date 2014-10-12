var path = require('path');
var Sequelize = require('sequelize');

var database = require('../config/database');

global.Model = Sequelize;

global.model = new Sequelize(
    database.database,
    database.username,
    database.password,
    database
);

require("fs").readdirSync("./models").forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require("../models/" + file);
});