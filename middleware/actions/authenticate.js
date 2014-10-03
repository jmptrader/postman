var crypto = require('crypto');

Action.register('auth', function(args) {
  var md5 = crypto.createHash('md5');
  md5.update(this.client.secret + this.auth_key);
  if (args.result !== md5.digest('hex')) {
    this.command('exit', {
      "msg": 'secret not match'
    });
    return
  }
  this.auth = true;
});