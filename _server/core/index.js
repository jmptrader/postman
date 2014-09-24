var path = require('path');
var msgpack = require('msgpack');

const commandPrefix = "DATA";
const commandSuffix = "END";

global.Action = {
  _actionMap: {},
  register: function(action, callback) {
    this._actionMap[action] = callback;
  },
  handle: function(command) {
    var action = command.split('|')[0];
    if (!this._actionMap.hasOwnProperty(action)) {
      return
    }
    var args = msgpack.unpack(command.substr(action.length));
    this._actionMap[action].call(cleartextStream, args);
  }
};

// load all actions
require("fs").readdirSync("./actions").forEach(function(file) {
  if (path.extname(file) !== 'js') return;
  require("./actions/" + file);
});

cleartextStream.init = function() {
  // TODO check client ip
  // send auth require command
  cleartextStream.reset();
};

cleartextStream.reset = function() {
  cleartextStream.hasPrefix = false;
  cleartextStream.replyStr = '';
};

// close remote client connection
cleartextStream.end = function() {
  cleartextStream.socket.end();
};

// send command to client
cleartextStream.command = function(action, args) {
  var commandStr = [action, msgpack.pack(args)].join('|');
  [commandPrefix, commandStr, commandSuffix].forEach(function(buf) {
    cleartextStream.write(buf);
  });
};

cleartextStream.addListener('data', function(body) {
  var data = body.toString().trim();
  console.log(data === commandPrefix, data == commandPrefix);
  if (data === commandPrefix) {
    cleartextStream.hasPrefix = true;
    return
  }
  if (cleartextStream.hasPrefix && data === commandSuffix) {
    Action.handle(cleartextStream.replyStr);
    cleartextStream.reset();
    cleartextStream.hasPrefix = false;
    return
  }
  cleartextStream.replyStr += data;
});

cleartextStream.addListener('close', function() {
  // TODO: warning should be raised to tell administrator: client closed.
  // all command in queue will resend after reconnect.
});