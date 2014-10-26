var fs = require('fs');
var tls = require('tls');
var sys = require('sys');
var port = require('../config/domain').middleware.port;

// all online sender hash.
global.senderMap = {};

// require all models
require('../api/core/models');
require('../api/core/mail-sync');
require('./core/command');

var options = {
    requestCert: true,
    key: fs.readFileSync('../config/pems/middleware/private-key.pem'),
    cert: fs.readFileSync('../config/pems/middleware/public-cert.pem')
};

// start server and load core
tls.createServer(options, function (cleartextStream) {
    sys.puts("TLS connection established: " + cleartextStream.remoteAddress);

    // load listeners
    global.cleartextStream = cleartextStream;

    // init connection info
    require('./core')();

    cleartextStream.setEncoding('utf8');
    cleartextStream.pipe(cleartextStream);
}).listen(port);