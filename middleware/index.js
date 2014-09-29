var fs = require('fs');
var tls = require('tls');
var sys = require('sys');

var options = {
  requestCert: true,
  key: fs.readFileSync('config/pems/private-key.pem'),
  cert: fs.readFileSync('config/pems/public-cert.pem')
}

// start server and load core
tls.createServer(options, function(cleartextStream) {
  sys.puts("TLS connection established: " + cleartextStream.remoteAddress);

  // load listeners
  global.cleartextStream = cleartextStream;
  require('./core');
  // init connection info
  cleartextStream.init && cleartextStream.init();

  cleartextStream.setEncoding('utf8');
  cleartextStream.pipe(cleartextStream);
}).listen(9002);