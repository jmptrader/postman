var path = require('path');
var sys = require('sys');
var crypto = require('crypto');

const commandPrefix = "DATA";
const commandSuffix = "END";

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
        this._actionMap[action].call(cleartextStream, args);
    }
};

// load all actions
require("fs").readdirSync("./actions").forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require("../actions/" + file);
});

// check client ip
var init = function () {
    var c = this;
    Client.find({
        where: {
            ip: c.remoteAddress
        }
    }).success(function (client) {
        if (!client) {
            c.command('exit', {
                "msg": 'client not found'
            })
            return;
        }
        // send auth require command
        crypto.randomBytes(16, function (ex, buf) {
            c.auth_key = buf.toString('hex').substr(0, 16);
            c.client = client;
            c.auth = false;
            c.command('helo', {
                auth_key: c.auth_key
            });
        });
    });
}

module.exports = function () {
    var c = cleartextStream;
    init.call(c);

    // send command to client
    c.command = function (action, args) {
        crypto.randomBytes(16, function (ex, buf) {
            var commandStr = ["+" + buf.toString('hex').substr(0, 16), action, JSON.stringify(args)].join('|');
            c.write(commandStr + '\n');
        });
    };

    // receive data and parse it
    c.addListener('data', function (data) {
        Action.handle(data.trim());
    });

    // trigger when client close
    c.addListener('close', function () {
        sys.puts("TLS connection closed");
        // TODO: warning should be raised to tell administrator: client closed.
        // all command in queue will resend after reconnect.
    });
};