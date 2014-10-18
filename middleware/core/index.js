var path = require('path');
var sys = require('sys');
var crypto = require('crypto');

var LINEFEED = '\f';

global.Action = {
    _actionMap: {},
    register: function (action, callback) {
        this._actionMap[action] = callback;
    },
    handle: function (command) {
        var cmdAtt = command.split('|');
        var id = cmdAtt[0];
        var action = cmdAtt[1];
        if (!this._actionMap.hasOwnProperty(action)) {
            return
        }
        var args = JSON.parse(command.substr(action.length + id.length + 2));
        this._actionMap[action].call(cleartextStream, args, id);
    }
};

// load all actions
require("fs").readdirSync("./actions").forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require("../actions/" + file);
});

// check sender ip
var init = function () {
    var c = this;
    Sender.find({
        where: {
            ip: c.remoteAddress,
            status: {
                ne: 'unverified'
            }
        }
    }).complete(function (err, sender) {
        if (err || !sender) {
            return c.command('exit', {
                "msg": 'sender not found'
            });
        }
        // send auth require command
        crypto.randomBytes(16, function (ex, buf) {
            c.auth_key = buf.toString('hex').substr(0, 16);
            c.sender = sender;
            c.auth = false;
            c.command('helo', {
                auth_key: c.auth_key
            });
        });
    });
};

module.exports = function () {
    var c = cleartextStream;
    init.call(c);
    // send command to sender
    c.command = function (action, args) {
        crypto.randomBytes(4, function (ex, buf) {
            var commandStr = ["+" + buf.toString('hex').substr(0, 4), action, JSON.stringify(args)].join('|');
            c.write(commandStr + LINEFEED);
        });
    };

    c.response = function (id, res) {
        c.command('response', {
            id: id,
            body: res
        })
    };

    c.unhandledBuffer = '';
    // receive data and parse it
    c.addListener('data', function (data) {
        data = c.unhandledBuffer + data;
        var buf = data.split(LINEFEED);
        c.unhandledBuffer = buf.pop();
        buf.forEach(function (cmd) {
            Action.handle(cmd.trim());
        });
    });
    // trigger when sender close
    c.addListener('close', function () {
        if (c.sender) {
            c.sender.status = 'offline';
            c.sender.save(['status']).success(function () {
                sys.puts('TLS connection ' + c.sender.ip + ' closed');
            });
        }
    });
};