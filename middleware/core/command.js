var path = require('path');

global.Command = {
    _commandMap: {},
    register: function (command, callback) {
        this._commandMap[command] = callback;
    },
    handle: function (sender_id, cmd) {
        if (!senderMap[sender_id]) {
            return setTimeout(function () {
                Command.handle(sender_id, cmd);
            }, 1000);
        }
        if (this._commandMap[cmd.command]) {
            var command = cmd.command;
            delete(cmd.command);
            this._commandMap[command].call(senderMap[sender_id], cmd);
        }
    }
};

require("fs").readdirSync(path.join(__dirname, '../commands')).forEach(function (file) {
    if (path.extname(file) !== '.js') return;
    require(path.join(__dirname, '../commands', file));
});

var commandLoop = function () {
    redisBlockClient.BRPOP('jianxin:command', 0, function (err, cmdArr) {
        if (err) return commandLoop();
        var sender_id = cmdArr[1].split(':')[0];
        var cmd = JSON.parse(cmdArr[1].substr(sender_id.length + 1));
        if (Command._commandMap[cmd.command]) {
            Command.handle(sender_id, cmd);
        }
        commandLoop();
    });
};

// start command queue listen.
commandLoop();