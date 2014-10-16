var crypto = require('crypto');
var express = require('express');
var path = require('path');

module.exports = global.router = express.Router({
    mergeParams: true
});

// sender auth
router.use(function (req, res, next) {
    if (!req.param('params')) {
        return res.jsonp({code: 422, msg: 'params needed.'});
    }
    var params;
    try {
        params = JSON.parse(req.param('params'));
    } catch (e) {
        return res.jsonp({code: 406, msg: 'params not acceptable.'});
    }
    if (!params['expire'] || params['expire'] < Date.now() / 1000) {
        return res.jsonp({code: 401, msg: 'secret has expired.'});
    }
    Sender.find({
        where: {
            domain: req.params['domain'],
            status: {
                ne: 'unverified'
            }
        }
    }).success(function (sender) {
        if (!sender) {
            return res.jsonp({code: 404, msg: 'sender not found.'});
        }
        var md5sum = crypto.createHash('md5');
        md5sum.update(req.param('params') + sender.api_key);
        if (req.params['secret'] !== md5sum.digest('hex')) {
            return res.jsonp({code: 401, msg: 'sender unauthorized.'});
        }
        req.sender = sender;
        next();
    });
});

require("fs").readdirSync("./actions").forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require("../actions/" + file);
});