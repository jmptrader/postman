var redis = require('redis');
var redisConfig = require('../../config/database').redis;

var redisClient = redis.createClient(
    redisConfig.port,
    redisConfig.host,
    redisConfig.options
);

global.newCommand = function (sender_id, command, next) {
    redisClient.LPUSH('jianxin:command', [sender_id, JSON.stringify(command)].join(':'), next)
};
