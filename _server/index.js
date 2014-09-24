var fs = require('fs');
var tls = require('tls');
var sys = require('sys');

configDir = process.env.POSTMAN_CONFIG_DIR;

var options = {
  requestCert: true,
  key: fs.readFileSync(configDir + '/ssl/private-key.pem'),
  cert: fs.readFileSync(configDir + '/ssl/public-cert.pem')
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
}).listen(8000);